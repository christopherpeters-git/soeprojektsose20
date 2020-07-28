package library

import (
	"database/sql"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"
)

//Sorts an array ascending alphabetical order
func sortArray(video []Video) {
	sort.Slice(video, func(i, j int) bool {
		return video[i].Show < video[j].Show
	})
}

//Check if string is legal
func IsStringLegal(str string) bool {
	for _, c := range str {
		if !strings.Contains(LetterBytes, strings.ToLower(string(c))) {
			return false
		}
	}
	return true
}

//Creates the array of favorite videos for a given user
func FillUserVideoArray(user *User, userDB *sql.DB) error {
	//Getting the informations about the user
	rows, err := userDB.Query("select users_username,video from user_has_favorite_videos where Users_Username = ?", user.Username)
	if err != nil {
		return errors.New("SQL query failed: \n" + err.Error())
	}
	user.FavoriteVideos = make([]Video, 0)
	var username string
	var videoStr string
	for rows.Next() {
		err := rows.Scan(&username, &videoStr)
		if err != nil {
			return errors.New("Scanning rows failed: \n" + err.Error())
		}
		var video Video
		log.Println(videoStr)
		err = json.Unmarshal([]byte(videoStr), &video)
		if err != nil {
			return errors.New("unmarshalling failed: \n" + err.Error())
		}
		user.FavoriteVideos = append(user.FavoriteVideos, video)
	}
	return nil
}

//Converts a Videoarray map to an video array
func ConvertMapToArray(mapToConvert map[string][]Video) []Video {
	var channelArray []Video
	var channelMap = mapToConvert
	for _, v := range channelMap {
		for video := range v {
			channelArray = append(channelArray, v[video])
		}
	}
	sortArray(channelArray)
	return channelArray
}

//Generate random session id
func GenerateSessionId(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = LetterBytes[rand.Intn(len(LetterBytes))]
	}
	return string(b)
}

//Sorts allVideos after Channels and Shows and loads them into a map
func SortByChannelAndShow(allVideos []Video) map[string]map[string][]Video { //channel          //show Array
	log.Println("Started sorting...")
	var videosSortedAfterChannels = make(map[string]map[string][]Video)
	for _, video := range allVideos {
		if videosSortedAfterChannels[video.Channel] == nil {
			videosSortedAfterChannels[video.Channel] = make(map[string][]Video)
		}
		if videosSortedAfterChannels[video.Channel][video.Show] == nil {
			videosSortedAfterChannels[video.Channel][video.Show] = make([]Video, 0)
		}
		videosSortedAfterChannels[video.Channel][video.Show] = append(videosSortedAfterChannels[video.Channel][video.Show], video)
	}
	log.Println("Sorting finished!")
	return videosSortedAfterChannels
}

//Starts a connection to a database and adds it to the connection map
func InitDataBaseConnection(dbConnections map[string]*sql.DB, driverName string, user string, password string, url string, dbName string, idName string) *sql.DB {
	//Open and check the sql-databse connection
	db, err := sql.Open(driverName,
		user+":"+password+"@tcp("+url+")/"+dbName)
	if err != nil {
		log.Fatal("Database connection failed: " + err.Error())
		return nil
	}
	dbConnections[idName] = db
	return db
}

//Reports an error and prints to the log
func ReportError(w http.ResponseWriter, statusCode int, responseMessage string, logMessage string) {
	http.Error(w, responseMessage, statusCode)
	log.Println(logMessage)
}

//Reports a DetailedError
func ReportDetailedError(w http.ResponseWriter, dErr DetailedHttpError) {
	http.Error(w, dErr.PublicError(), dErr.Status())
	log.Println(dErr.Error())
}

//Loads the good.json and stores it on the server
func ParseVideosFromJson(videos *[]Video) {
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
		*videos = append(*videos, video)
	}
}

//Checks if a request has a valid cookie and writes to user if valid cookies exists
func IsUserLoggedInWithACookie(r *http.Request, userDB *sql.DB, user *User) DetailedHttpError {
	//Check if there is a cookie in the request
	cookie, err := r.Cookie(AuthCookieName)
	if err == nil {
		if cookie.Value == "0" {
			return &ClientError{http.StatusForbidden, AuthenticationFailedResponse, errors.New("session ID in cookie is 0")}
		}
		rows, err := userDB.Query("select id,name,username,passwordhash,session_id from users where session_id = ?", cookie.Value)
		if err != nil {
			return &ServerError{http.StatusInternalServerError, InternalServerErrorResponse, errors.New("SQL query failed: \n" + err.Error())}
		}
		if rows.Next() {
			err := rows.Scan(&user.Id, &user.Name, &user.Username, &user.passwordHash, &user.sessionId)
			if err != nil {
				return &ServerError{http.StatusInternalServerError, InternalServerErrorResponse, errors.New("Scanning rows failed: \n" + err.Error())}
			}
		} else {
			log.Println("No such SessionId found")
			return &ClientError{http.StatusForbidden, AuthenticationFailedResponse, errors.New("No user found with session id: " + cookie.Value)}
		}
		return nil
	} else {
		return &ClientError{http.StatusForbidden, AuthenticationFailedResponse, errors.New("no cookie found with name: " + AuthCookieName)}
	}
}

//Returns true if user exists and fills user with user information from the DB
func LoginUser(userDB *sql.DB, user *User, incomingUsername string, incomingPassword string) DetailedHttpError {
	//Get userdata from db for comparison
	rows, err := userDB.Query("select id,name,username,passwordhash,session_id from users where username = ?", incomingUsername)
	if err != nil {
		return &ServerError{http.StatusInternalServerError, InternalServerErrorResponse, errors.New("sql query failed: \n" + err.Error())}
	}
	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Name, &user.Username, &user.passwordHash, &user.sessionId)
		if err != nil {
			return &ServerError{http.StatusInternalServerError, InternalServerErrorResponse, errors.New("scanning rows failed: \n" + err.Error())}
		}
		//Compare found password hash with incoming password hash
		err = bcrypt.CompareHashAndPassword([]byte(user.passwordHash), []byte(incomingPassword))
		if err != nil {
			return &ClientError{http.StatusForbidden, "wrong password", errors.New("wrong password: \n" + err.Error())}
		} else {
			log.Println("Entered Password is correct")
		}
	} else {
		return &ClientError{http.StatusNotFound, "user not found", errors.New("empty sql result set")}
	}
	return nil
}

//Places a cookie on a client (w) and stores the session id for the given username
func PlaceCookie(w http.ResponseWriter, userDB *sql.DB, incomingUsername string) DetailedHttpError {
	//Generating and inserting a new SessionId in the db
	sessionId := GenerateSessionId(255)
	_, err := userDB.Exec("UPDATE users set Session_Id = ? where username = ?", sessionId, incomingUsername)
	if err != nil {
		return &ServerError{http.StatusInternalServerError, InternalServerErrorResponse, errors.New("Updating sql failed: \n" + err.Error())}
	}
	expire := time.Now().AddDate(0, 0, 2)
	cookie := http.Cookie{
		Name:       AuthCookieName,
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
	return nil
}
