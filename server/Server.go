package main

//TODO
/*
-Server gibt nur Videos von channel wieder, welcher in der request übergebenwird
*/

import (
	lib "./library"
	"database/sql"
	"encoding/json"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

//**********************<Constants>**********************************
//URL's
const IncomingGetVideosRequestUrl = "/getVideos/"
const IncomingPostUserRequestUrl = "/login/"
const IncomingPostRegisterRequestUrl = "/register/"
const IncomingGetVideosFromChannelRequestUrl = "/getVideoByChannel"
const IncomingGetVideoClickedRequestUrl = "/clickVideo"
const IncomingPostAddToFavoritesRequestUrl = "/addToFavorites/"
const IncomingPostLogoutRequestUrl = "/logout/"

//Parameter
const ChannelNameParameter = "channel"
const VideoTitleParameter = "videoTitle"

//Error messages
const InternalServerErrorResponse = "Internal server error - see logs"

//db connection names
const UserDBconnectionName = "userdb"

//Paths
const CrawlerDirName = "crawler"
const VideoJsonPath = CrawlerDirName + "/good.json"

//Characters allowed in SessionID
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const authCookieName = "mediathekauth"

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
	sessionId      string
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
	sessionId    string
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
	err = lib.DownloadJson(CrawlerDirName, VideoJsonPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	parseVideosFromJson()
	sortByChannelAndShow()
	//Connect do database
	defer initDataBaseConnection("mysql", "root", "soe2020", "localhost:3306", UserDBconnectionName, UserDBconnectionName).Close()

	log.Print("Server has started...")
	http.Handle("/", http.FileServer(http.Dir("test_frontend/")))
	http.HandleFunc(IncomingGetVideosRequestUrl, handleGetAllVideos)
	http.HandleFunc(IncomingGetVideosFromChannelRequestUrl, handleGetVideosFromChannel)
	http.HandleFunc(IncomingPostUserRequestUrl, handlePostLogin)
	http.HandleFunc(IncomingPostAddToFavoritesRequestUrl, handlePostAddVideoToFavorites)
	http.HandleFunc(IncomingGetVideoClickedRequestUrl, handleGetVideoClicked)
	http.HandleFunc(IncomingPostRegisterRequestUrl, handlePostRegisterUser)
	http.HandleFunc(IncomingPostLogoutRequestUrl, handlePostLogout)
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Starting Server failed: " + err.Error())
	}
}

//**************************************<Helpers>************************************************************************
//Generate random session id
func generateSessionId(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

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
	byteValue, err := ioutil.ReadFile(VideoJsonPath)
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
}

//Checks if a request has a valid cookie and writes to user if valid cookies exists
func isUserLoggedInWithACookie(r *http.Request, userDB *sql.DB, user *User) (bool, error) {
	//Check if there is a cookie in the request
	cookie, err := r.Cookie(authCookieName)
	if err == nil {
		if cookie.Value == "0" {
			return false, nil
		}
		rows, err := userDB.Query("Select * from Users where Session_Id = ?", cookie.Value)
		if err != nil {
			return false, errors.New("SQL query failed: \n" + err.Error())
		}
		if rows.Next() {
			err := rows.Scan(&user.Id, &user.Name, &user.Username, &user.passwordHash, &user.sessionId)
			if err != nil {
				return false, errors.New("Scanning rows failed: \n" + err.Error())
			}
		} else {
			log.Println("No such SessionId found")
			return false, nil
		}
		return true, nil
	} else {
		log.Println("No cookie found with name " + authCookieName)
		return false, nil
	}
}

func loginUser(w http.ResponseWriter, userDB *sql.DB, user *User, incomingUsername string, incomingPassword string) bool {
	//Get userdata from db for comparison
	rows, err := userDB.Query("select * from users where username = ?", incomingUsername)
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "Sql query failed: \n"+err.Error())
		return false
	}
	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Name, &user.Username, &user.passwordHash, &user.sessionId)
		if err != nil {
			reportError(w, 500, InternalServerErrorResponse, "Scanning rows failed: \n"+err.Error())
			return false
		}
		log.Println("User from Query: username: " + user.Username + " passwordhash: " + user.passwordHash)
		//Compare found password hash with incoming password hash
		err = bcrypt.CompareHashAndPassword([]byte(user.passwordHash), []byte(incomingPassword))
		if err != nil {
			reportError(w, 400, "Wrong password", "Wrong password: \n"+err.Error())
			return false
		} else {
			log.Println("Entered Password is correct")
		}
		//Generating and inserting a new SessionId in the db
		sessionId := generateSessionId(255)
		_, err = userDB.Exec("UPDATE users set Session_Id = ? where username = ?", sessionId, incomingUsername)
		if err != nil {
			reportError(w, 500, InternalServerErrorResponse, "Updating sql failed: \n"+err.Error())
			return false
		}
		expire := time.Now().AddDate(0, 0, 2)
		cookie := http.Cookie{
			Name:       authCookieName,
			Value:      sessionId,
			Path:       "/",
			Domain:     "localhost",
			Expires:    expire,
			RawExpires: expire.Format(time.UnixDate),
			MaxAge:     172800,
			Secure:     false,
			HttpOnly:   true,
			SameSite:   http.SameSiteLaxMode,
		}
		http.SetCookie(w, &cookie)
	} else {
		reportError(w, 404, "User not found", "Empty sql result set \n")
		return false
	}
	return true
}

//**************************************</Helpers>************************************************************************

//**************************************<Handlers>***********************************************************************
func handlePostLogout(w http.ResponseWriter, r *http.Request) {
	log.Println("Answering handlePostLogout request started...")
	userDB := dbConnections[UserDBconnectionName]
	err := userDB.Ping()
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "Database connection failed: \n"+err.Error())
		return
	}
	var user User
	isCookieValid, err := isUserLoggedInWithACookie(r, userDB, &user)
	if !isCookieValid {
		if err != nil {
			reportError(w, 401, "Authentication failed", "Database connection failed: \n"+err.Error())
			return
		}
	}
	_, err = userDB.Exec("UPDATE users SET session_id = '0' WHERE username = ?", user.Username)
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "SQL update failed: \n"+err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("Logged out"))
	log.Println("Answered handlePostLogout successfully")
}

func handleGetVideoClicked(w http.ResponseWriter, r *http.Request) {
	log.Println("Answering handleGetVideoClicked request started...")
	viewcount := 0
	//Checking db connection
	userDB := dbConnections[UserDBconnectionName]
	err := userDB.Ping()
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "Database connection failed: \n"+err.Error())
		return
	}
	queryResults, ok := r.URL.Query()[VideoTitleParameter]
	if !ok || len(queryResults) < 1 {
		reportError(w, 400, "url parameter unknown", "Cant find parameter "+VideoTitleParameter)
		return
	}
	title := queryResults[0]
	log.Println("Content of parameter '" + ChannelNameParameter + "': " + title)
	//Get
	rows, err := userDB.Query("select * from videos where VideoTitle = ?", title)
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "Sql query failed: \n"+err.Error())
		return
	}
	if rows.Next() {
		var tempTitle string
		var tempViews int
		err = rows.Scan(&tempTitle, &tempViews)
		if err != nil {
			reportError(w, 500, InternalServerErrorResponse, "Scanning rows failed: \n"+err.Error())
			return
		}
		tempViews++
		viewcount = tempViews
		_, err = userDB.Exec("UPDATE videos SET Views = ? WHERE VideoTitle = ?", tempViews, tempTitle)
		if err != nil {
			reportError(w, 500, InternalServerErrorResponse, "SQL update failed: \n"+err.Error())
			return
		}
	} else {
		_, err = userDB.Exec("INSERT INTO videos (VideoTitle,Views) \n Values(?,?) ", title, 1)
		if err != nil {
			viewcount = 1
			reportError(w, 500, InternalServerErrorResponse, "SQL update failed: \n"+err.Error())
			return
		}
	}
	viewCountStr := strconv.Itoa(viewcount)
	w.WriteHeader(200)
	w.Write([]byte(viewCountStr))
	log.Println("Answered handleGetVideoClicked successfully")
}

func handlePostAddVideoToFavorites(w http.ResponseWriter, r *http.Request) {
	log.Println("Answering handlePostAddVideoToFavorites request started...")
	//Checking db connection
	userDB := dbConnections[UserDBconnectionName]
	err := userDB.Ping()
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "Database connection failed: \n"+err.Error())
		return
	}
	incomingVideo := r.FormValue("video")
	var user User
	isCookieValid, err := isUserLoggedInWithACookie(r, userDB, &user)
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "Cookie validation failed: \n"+err.Error())
		return
	}
	if !isCookieValid {
		w.WriteHeader(401)
		w.Write([]byte("Access denied"))
		return
	}
	//Check if video is already in favorites
	rows, err := userDB.Query("Select * from user_has_favorite_videos where Users_Username = ? and Video = ?", user.Username, incomingVideo)
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "SQL query checking for username and video failed: \n"+err.Error())
		return
	}
	if rows.Next() {
		reportError(w, 400, "Video is already favorite", "Video is already in favorites: \n")
		return
	}
	//Insert new video for the user
	_, err = userDB.Exec("INSERT INTO user_has_favorite_videos (Users_Username,Video) \n Values(?,?)", user.Username, incomingVideo)
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "SQL insert failed: \n"+err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("Added video to favorites successfully"))
	log.Println("Answering handlePostAddVideoToFavorites request started...")
}

func handleGetVideosFromChannel(w http.ResponseWriter, r *http.Request) {
	log.Println("Answering handleGetAllVideos request started...")
	queryResults, ok := r.URL.Query()[ChannelNameParameter]
	if !ok || len(queryResults) < 1 {
		reportError(w, 400, "url parameter unknown", "Cant find parameter "+ChannelNameParameter)
		return
	}
	channel := queryResults[0]
	log.Println("Content of parameter '" + ChannelNameParameter + "': " + channel)
	resultSetInBytes, err := json.MarshalIndent(videosSortedAfterChannels[channel], "", " ")
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "Marshaling failed: \n"+err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write(resultSetInBytes)
	log.Println("Answered handleGetAllVideos request successfully...")
}

func handlePostRegisterUser(w http.ResponseWriter, r *http.Request) {
	log.Print("answering handlePostRegisterUser request ...")

	//Checking db connection
	userDB := dbConnections[UserDBconnectionName]
	err := userDB.Ping()
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "Database connection failed: \n"+err.Error())
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
		reportError(w, 500, InternalServerErrorResponse, "Sql query failed: \n"+err.Error())
		return
	}
	//Check if username is in use
	var user DB_User
	for rows.Next() {
		err = rows.Scan(&user.id, &user.name, &user.username, &user.passwordHash, &user.sessionId)
		if err != nil {
			reportError(w, 500, InternalServerErrorResponse, "Scanning rows failed: \n"+err.Error())
			return
		}
		if user.username == incomingUsername {
			reportError(w, 400, "User already exists", "Username taken: "+user.username)
			return
		}
	}
	//Hash incoming password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(incomingPassword), bcrypt.MinCost)
	log.Printf("User created: Name: %s username: %s passwordhash: %s", incomingName, incomingUsername, string(passwordHash))
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "Hashing password '"+incomingPassword+"' failed: \n"+err.Error())
		return
	}
	//Create user in database
	_, err = userDB.Exec("INSERT INTO users (Name,Username,PasswordHash)\nValues(?,?,?)", incomingName, incomingUsername, string(passwordHash))
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "SQL insert failed: \n"+err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("Created a new User account"))
	log.Println("answered handlePostRegisterUser request successfully")
}

func handleGetAllVideos(w http.ResponseWriter, r *http.Request) {
	log.Print("Answering handleGetAllVideos request...")
	//Writing the result set to the responseWriter as a json-string
	resultSetInBytes, err := json.MarshalIndent(allVideos, "", " ")
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "Marshaling failed: \n"+err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write(resultSetInBytes)
	log.Print("Answered handleGetAllVideos request successfully...")
}

func handlePostLogin(w http.ResponseWriter, r *http.Request) {
	log.Print("answering handlePostLogin request ...")

	//Checking db connection
	userDB := dbConnections[UserDBconnectionName]
	var user User
	err := userDB.Ping()
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "Database connection failed: \n"+err.Error())
		return
	}
	validCookieFound, err := isUserLoggedInWithACookie(r, userDB, &user)
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "Cookie validation failed: \n"+err.Error())
		return
	}

	log.Println("No cookie found with name: " + authCookieName)
	if !validCookieFound {
		//Parse username and password from request
		err = r.ParseForm()
		if err != nil {
			reportError(w, 400, "Invalid request parameters", "Parameter parsing error: "+err.Error())
		}
		incomingUsername := r.FormValue("usernameInput")
		incomingPassword := r.FormValue("passwordInput")

		if !loginUser(w, userDB, &user, incomingUsername, incomingPassword) {
			return
		}
	}
	//Getting the informations about the user
	rows, err := userDB.Query("select * from user_has_favorite_videos where Users_Username = ?", user.Username)
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "Sql query failed: \n"+err.Error())
		return
	}
	user.FavoriteVideos = make([]Video, 0)
	var username string
	var videoStr string
	for rows.Next() {
		err := rows.Scan(&username, &videoStr)
		if err != nil {
			reportError(w, 500, InternalServerErrorResponse, "Scanning rows failed: \n"+err.Error())
			return
		}
		var video Video
		err = json.Unmarshal([]byte(videoStr), &video)
		if err != nil {
			reportError(w, 500, InternalServerErrorResponse, "unmarshalling failed: \n"+err.Error())
			return
		}
		user.FavoriteVideos = append(user.FavoriteVideos, video)
	}

	//Writing the result set to the responseWriter as a json-string
	userInBytes, err := json.MarshalIndent(user, "", "   ")
	if err != nil {
		reportError(w, 500, InternalServerErrorResponse, "Marshaling failed: \n"+err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write(userInBytes)
	log.Print("Answered handlePostLogin request successfully...")
}

//**************************************</Handlers>***********************************************************************
