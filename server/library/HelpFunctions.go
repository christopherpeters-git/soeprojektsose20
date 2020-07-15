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
	"time"
)

func ConvertMapToArray(mapToConvert map[string][]Video) []Video {
	var channelArray []Video
	var channelMap = mapToConvert
	for _, v := range channelMap {
		for video := range v {
			channelArray = append(channelArray, v[video])
		}
	}
	return channelArray
}

//Generate random session id
func generateSessionId(n int) string {
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

func InitDataBaseConnection(dbConnections map[string]*sql.DB, driverName string, user string, password string, url string, dbName string, idName string) *sql.DB {
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

func ReportError(w http.ResponseWriter, statusCode int, responseMessage string, logMessage string) {
	http.Error(w, responseMessage, statusCode)
	log.Println(logMessage)
}

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
func IsUserLoggedInWithACookie(r *http.Request, userDB *sql.DB, user *User) (bool, error) {
	//Check if there is a cookie in the request
	cookie, err := r.Cookie(AuthCookieName)
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
		log.Println("No cookie found with name " + AuthCookieName)
		return false, nil
	}
}

func LoginUser(w http.ResponseWriter, userDB *sql.DB, user *User, incomingUsername string, incomingPassword string) bool {
	//Get userdata from db for comparison
	rows, err := userDB.Query("select * from users where username = ?", incomingUsername)
	if err != nil {
		ReportError(w, 500, InternalServerErrorResponse, "Sql query failed: \n"+err.Error())
		return false
	}
	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Name, &user.Username, &user.passwordHash, &user.sessionId)
		if err != nil {
			ReportError(w, 500, InternalServerErrorResponse, "Scanning rows failed: \n"+err.Error())
			return false
		}
		log.Println("User from Query: username: " + user.Username + " passwordhash: " + user.passwordHash)
		//Compare found password hash with incoming password hash
		err = bcrypt.CompareHashAndPassword([]byte(user.passwordHash), []byte(incomingPassword))
		if err != nil {
			ReportError(w, 400, "Wrong password", "Wrong password: \n"+err.Error())
			return false
		} else {
			log.Println("Entered Password is correct")
		}
		//Generating and inserting a new SessionId in the db
		sessionId := generateSessionId(255)
		_, err = userDB.Exec("UPDATE users set Session_Id = ? where username = ?", sessionId, incomingUsername)
		if err != nil {
			ReportError(w, 500, InternalServerErrorResponse, "Updating sql failed: \n"+err.Error())
			return false
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
	} else {
		ReportError(w, 404, "User not found", "Empty sql result set \n")
		return false
	}
	return true
}
