package main

import (
	lib "./library"
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//**********************<Constants>**********************************
//URL's
var IncomingGetVideosRequestUrl = "/getVideos/"
var IncomingPostUserRequestUrl = "/login/"

//Error messages
var InternalServerErrorResponse = "Internal server error - see logs"

//db connection names
var userDBconnectionName = "videoDB"

//Paths
var crawlerDirName = "crawler"
var videoJsonPath = crawlerDirName + "/good.json"

//**********************</Constants>**********************************

type Video struct {
	channel     string
	title       string
	show        string
	releaseDate string
	duration    string
	link        string
	pageLink    string
	fileName    string //Shouldnt be used
}

//User credentials
type UserEntry struct {
	username     string
	passwordHash string
	//sessionId    uint64
}

//User data (favorites, ...)
type UserData struct {
	username       string
	name           string
	favoriteVideos string
}

var videos = make([]Video, 0)

var dbConnections = make(map[string]*sql.DB)

func main() {
	//Creates a log file
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	f, err := os.OpenFile("Server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	//Start crawler
	err = lib.DownloadJson(crawlerDirName, videoJsonPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	parseVideosFromJson()
	//Connect do database
	defer initDataBaseConnection("mysql", "root", "soe2020", "localhost:3306", userDBconnectionName, userDBconnectionName).Close()
	log.Print("Server has started...")

	http.Handle("/", http.FileServer(http.Dir("test_frontend/")))
	http.HandleFunc(IncomingGetVideosRequestUrl, handleGetVideos)
	http.HandleFunc(IncomingPostUserRequestUrl, handleGetUserInformation)
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Starting Server failed: " + err.Error())
	}
}

//**************************************<Helpers>************************************************************************
func initDataBaseConnection(driverName string, user string, password string, url string, dbName string, idName string) *sql.DB {
	//Open and check the sql-databse connection
	db, err := sql.Open(driverName,
		user+":"+password+"@tcp("+url+")/"+dbName)
	if err != nil {
		log.Fatal("Database connection failed: " + err.Error())
		return nil
	}
	//defer db.Close() //???????
	dbConnections[idName] = db
	return db
}

func reportError(w http.ResponseWriter, statusCode int, responseMessage string, logMessage string) {
	http.Error(w, responseMessage, statusCode)
	log.Println(logMessage)
}

func parseVideosFromJson() {
	data := make([][]string, 0)
	byteValue, err := ioutil.ReadFile(videoJsonPath)
	if err != nil {
		log.Println("Parsing failed: " + err.Error())
		return
	}

	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		log.Println("Parsing failed: " + err.Error())
		return
	}

	for _, videoEntry := range data {
		video := Video{
			channel:     videoEntry[0],
			title:       videoEntry[1],
			show:        videoEntry[2],
			releaseDate: videoEntry[3],
			duration:    videoEntry[4],
			link:        videoEntry[5],
			pageLink:    videoEntry[6],
			fileName:    videoEntry[7],
		}
		videos = append(videos, video)
	}

	//Testing
	//for _, a := range videos{
	//	log.Print(a.channel)
	//	log.Print(a.duration)
	//	log.Print(a.fileName)
	//	log.Print(a.link)
	//	log.Print(a.pageLink)
	//	log.Print(a.releaseDate)
	//	log.Print(a.show)
	//	log.Print(a.title)
	//	log.Println()
	//}
}

//**************************************</Helpers>************************************************************************

//**************************************<Handlers>***********************************************************************
func handleGetVideos(w http.ResponseWriter, r *http.Request) {
	log.Print("Answering handleGetVideos request...")
	//Writing the result set to the responseWriter as a json-string
	resultSetInBytes, err := json.MarshalIndent(videos, "", "   ")
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "Marshaling failed: \n"+err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write(resultSetInBytes)
	log.Print("Answered handleGetVideos request successfully...")
}

func handleGetUserCookie(w http.ResponseWriter, r *http.Request) {
	//TODO
	//Post requst with user credentials
}

func handleGetUserInformation(w http.ResponseWriter, r *http.Request) {
	//Unter der Annahme, dass kein login cookie verwendet wird

	log.Print("answering handleGetUserInformation request ...")

	//Checking db connection
	userDB := dbConnections[userDBconnectionName]
	err := userDB.Ping()
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "Database connection failed: \n"+err.Error())
		return
	}
	//Parse username and password from request
	err = r.ParseForm()
	if err != nil {
		reportError(w, 400, "Invalid request parameters", "Parameter parsing error: "+err.Error())

	}
	incomingUsername := r.FormValue("username")
	incomingPassword := r.FormValue("password")
	//Get userdata from db for comparison
	rows, err := userDB.Query("select * from user where username = ?", incomingUsername)
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "Sql query failed: \n"+err.Error())
		return
	}
	var userCredentials UserEntry
	if rows.Next() {
		err := rows.Scan(&userCredentials.username, &userCredentials.passwordHash)
		if err != nil {
			reportError(w, 500, InternalServerErrorResponse, "Scanning rows failed: \n"+err.Error())
			return
		}
		log.Println("User from Query: username: " + userCredentials.username + " name: " + userCredentials.passwordHash)
		//Compare found password hash with incoming password hash
		err = bcrypt.CompareHashAndPassword([]byte(userCredentials.passwordHash), []byte(incomingPassword))
		if err != nil {
			reportError(w, 400, "Wrong password", "Wrong password: \n"+err.Error())
			return
		} else {
			log.Println("Entered Password is correct")
		}
	} else {
		reportError(w, 404, "Wrong username", "Empty sql result set: \n"+err.Error())
		return
	}

	//Getting the informations about the user
	rows, err = userDB.Query("select * from userinformations where username = ?", incomingUsername)
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "Sql query failed: \n"+err.Error())
		return
	}
	var userInformation UserData
	if rows.Next() {
		err := rows.Scan(&userInformation.username, &userInformation.name, &userInformation.favoriteVideos)
		if err != nil {
			reportError(w, 500, InternalServerErrorResponse, "Scanning rows failed: \n"+err.Error())
			return
		}
	} else {
		reportError(w, 500, InternalServerErrorResponse, "Empty sql result set: \n"+err.Error())
		return
	}

	//Writing the result set to the responseWriter as a json-string
	userinformationInBytes, err := json.MarshalIndent(userInformation, "", "   ")
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "Marshaling failed: \n"+err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write(userinformationInBytes)
	log.Print("Answered handleGetUserInformation request successfully...")
}

//**************************************</Handlers>***********************************************************************
