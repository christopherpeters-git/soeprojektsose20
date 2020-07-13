package main

//TODO
/*
-Server gibt nur Videos von channel wieder, welcher in der request übergebenwird
*/

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
const IncomingGetVideosRequestUrl = "/getVideos/"
const IncomingPostUserRequestUrl = "/login/"
const IncomingPostRegisterRequestUrl = "/register/"
const IncomingGetVideosFromChannelRequestUrl = "/getVideoByChannel"

//Parameter
const channelNameParameter = "channel"

//Error messages
const internalServerErrorResponse = "Internal server error - see logs"

//db connection names
const userDBconnectionName = "userdb"

//Paths
const crawlerDirName = "crawler"
const videoJsonPath = crawlerDirName + "/good.json"

//**********************</Constants>**********************************

type Video struct {
	Channel     string `json:"channel"`
	Title       string `json:"title"`
	Show        string `json:"show"`
	ReleaseDate string `json:"releaseDate"`
	Duration    string `json:"duration"`
	Link        string `json:"link"`
	PageLink    string `json:"pageLink"`
	FileName    string `json:"fileName"`
}

type User struct {
	Id             string
	Name           string
	Username       string
	passwordHash   string
	FavoriteVideos []Video
}

type DB_Video struct {
	videoLink string
	views     int
}

type DB_User struct {
	id           string
	name         string
	username     string
	passwordHash string
}

var allVideos = make([]Video, 0)
var videosSortedAfterChannels = make(map[string]map[string][]Video)

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
	sortByChannelAndShow()
	//Connect do database
	defer initDataBaseConnection("mysql", "root", "soe2020", "localhost:3306", userDBconnectionName, userDBconnectionName).Close()

	log.Print("Server has started...")
	http.Handle("/", http.FileServer(http.Dir("test_frontend/")))
	http.HandleFunc(IncomingGetVideosRequestUrl, handleGetAllVideos)
	http.HandleFunc(IncomingGetVideosFromChannelRequestUrl, handleGetVideosFromChannel)
	http.HandleFunc(IncomingPostUserRequestUrl, handleGetUserInformation)
	http.HandleFunc(IncomingPostRegisterRequestUrl, handleRegisterUser)
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Starting Server failed: " + err.Error())
	}
}

//**************************************<Helpers>************************************************************************
//Sorts allVideos after Channels and Shows and loads them into a map
func sortByChannelAndShow() { //channel          //show Array
	log.Println("Started sorting...")
	for _, video := range allVideos {
		if videosSortedAfterChannels[video.Channel] == nil {
			videosSortedAfterChannels[video.Channel] = make(map[string][]Video)
		}
		if videosSortedAfterChannels[video.Channel][video.Show] == nil {
			videosSortedAfterChannels[video.Channel][video.Show] = make([]Video, 0)
		}
		videosSortedAfterChannels[video.Channel][video.Show] = append(videosSortedAfterChannels[video.Channel][video.Show], video)
	}
	//for k,v := range videosSortedAfterChannels{
	//   log.Println("key: "+ k)
	//   for k1, v1 := range v{
	//       log.Print(" show"+ k1)
	//       for _,v2 := range v1{
	//           log.Print(" "+ v2.Title)
	//       }
	//   }
	//}
	log.Println("Sorting finished!")
}

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
			Channel:     videoEntry[0],
			Title:       videoEntry[1],
			Show:        videoEntry[2],
			ReleaseDate: videoEntry[3],
			Duration:    videoEntry[4],
			Link:        videoEntry[5],
			PageLink:    videoEntry[6],
			FileName:    videoEntry[7],
		}
		allVideos = append(allVideos, video)
	}

	//Testing
	//for _, a := range allVideos{
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
func handleGetVideosFromChannel(w http.ResponseWriter, r *http.Request) {
	log.Print("Answering handleGetAllVideos request started...")
	queryResults, ok := r.URL.Query()[channelNameParameter]
	if !ok || len(queryResults) < 1 {
		reportError(w, 400, "url parameter unkown", "Cant find parameter "+channelNameParameter)
		return
	}
	channel := queryResults[0]
	log.Println("Content of parameter '" + channelNameParameter + "': " + channel)
	resultSetInBytes, err := json.MarshalIndent(videosSortedAfterChannels[channel], "", " ")
	if err != nil {
		reportError(w, 500, internalServerErrorResponse, "Marshaling failed: \n"+err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write(resultSetInBytes)
	log.Print("Answered handleGetAllVideos request successfully...")
}

func handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	log.Print("answering handleRegisterUser request ...")

	//Checking db connection
	userDB := dbConnections[userDBconnectionName]
	err := userDB.Ping()
	if err != nil {
		reportError(w, 500, internalServerErrorResponse, "Database connection failed: \n"+err.Error())
		return
	}
	//Parse username and password from request
	err = r.ParseForm()
	if err != nil {
		reportError(w, 400, "Invalid request parameters", "Parameter parsing error: "+err.Error())
		return
	}
	incomingUsername := r.FormValue("usernameInput")
	incomingName := r.FormValue("nameInput")
	incomingPassword := r.FormValue("passwordInput")
	//Check if any of the recieved information is empty
	if len(incomingUsername) < 1 || len(incomingName) < 1 || len(incomingPassword) < 1 {
		reportError(w, 400, "Send Information must not be empty", "one or more received strings is empty\n")
		return
	}
	//Get userdata from db for comparison
	rows, err := userDB.Query("select * from users where username = ?", incomingUsername)
	if err != nil {
		reportError(w, 500, internalServerErrorResponse, "Sql query failed: \n"+err.Error())
		return
	}
	//Check if username is in use
	var user DB_User
	for rows.Next() {
		err = rows.Scan(&user.id, &user.name, &user.username, &user.passwordHash)
		if err != nil {
			reportError(w, 500, internalServerErrorResponse, "Scanning rows failed: \n"+err.Error())
			return
		}
		if user.username == incomingUsername {
			reportError(w, 400, "Username taken", "Username taken: "+user.username)
			return
		}
	}
	//Hash incoming password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(incomingPassword), bcrypt.MinCost)
	log.Printf("User created: Name: %s username: %s passwordhash: %s", incomingName, incomingUsername, string(passwordHash))
	if err != nil {
		reportError(w, 500, internalServerErrorResponse, "Hashing password '"+incomingPassword+"' failed: \n"+err.Error())
		return
	}
	//Create user in database
	_, err = userDB.Exec("INSERT INTO users (Name,Username,PasswordHash)\nValues(?,?,?)", incomingName, incomingUsername, string(passwordHash))
	if err != nil {
		reportError(w, 500, internalServerErrorResponse, "SQL insert failed: \n"+err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("Created a new User account"))
	log.Println("answered handleRegisterUser request successfully")
}

func handleGetAllVideos(w http.ResponseWriter, r *http.Request) {
	log.Print("Answering handleGetAllVideos request...")
	//Writing the result set to the responseWriter as a json-string
	resultSetInBytes, err := json.MarshalIndent(allVideos, "", " ")
	if err != nil {
		reportError(w, 500, internalServerErrorResponse, "Marshaling failed: \n"+err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write(resultSetInBytes)
	log.Print("Answered handleGetAllVideos request successfully...")
}

func handleGetUserInformation(w http.ResponseWriter, r *http.Request) {
	log.Print("answering handleGetUserInformation request ...")

	//Checking db connection
	userDB := dbConnections[userDBconnectionName]
	err := userDB.Ping()
	if err != nil {
		reportError(w, 500, internalServerErrorResponse, "Database connection failed: \n"+err.Error())
		return
	}
	//Parse username and password from request
	err = r.ParseForm()
	if err != nil {
		reportError(w, 400, "Invalid request parameters", "Parameter parsing error: "+err.Error())
	}
	incomingUsername := r.FormValue("usernameInput")
	incomingPassword := r.FormValue("passwordInput")
	//Get userdata from db for comparison
	rows, err := userDB.Query("select * from users where username = ?", incomingUsername)
	if err != nil {
		reportError(w, 500, internalServerErrorResponse, "Sql query failed: \n"+err.Error())
		return
	}
	var user User
	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Name, &user.Username, &user.passwordHash)
		if err != nil {
			reportError(w, 500, internalServerErrorResponse, "Scanning rows failed: \n"+err.Error())
			return
		}
		log.Println("User from Query: username: " + user.Username + " passwordhash: " + user.passwordHash)
		//Compare found password hash with incoming password hash
		err = bcrypt.CompareHashAndPassword([]byte(user.passwordHash), []byte(incomingPassword))
		if err != nil {
			reportError(w, 400, "Wrong password", "Wrong password: \n"+err.Error())
			return
		} else {
			log.Println("Entered Password is correct")
		}
	} else {
		reportError(w, 404, "User not found", "Empty sql result set \n")
		return
	}

	//Getting the informations about the user
	rows, err = userDB.Query("select * from user_has_favorite_videos where Users_Username = ?", incomingUsername)
	if err != nil {
		reportError(w, 500, internalServerErrorResponse, "Sql query failed: \n"+err.Error())
		return
	}
	user.FavoriteVideos = make([]Video, 0)
	var username string
	var videoStr string
	for rows.Next() {
		err := rows.Scan(&username, &videoStr)
		if err != nil {
			reportError(w, 500, internalServerErrorResponse, "Scanning rows failed: \n"+err.Error())
			return
		}
		var video Video
		err = json.Unmarshal([]byte(videoStr), &video)
		if err != nil {
			reportError(w, 500, internalServerErrorResponse, "Scanning rows failed: \n"+err.Error())
			return
		}
		user.FavoriteVideos = append(user.FavoriteVideos, video)
	}

	//Writing the result set to the responseWriter as a json-string
	userInBytes, err := json.MarshalIndent(user, "", "   ")
	if err != nil {
		reportError(w, 500, internalServerErrorResponse, "Marshaling failed: \n"+err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write(userInBytes)
	log.Print("Answered handleGetUserInformation request successfully...")
}

//**************************************</Handlers>***********************************************************************
