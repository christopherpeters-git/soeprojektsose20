package main

import (
	lib "./library"
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var allVideos = make([]lib.Video, 0)
var videosSortedAfterChannels map[string]map[string][]lib.Video
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
	err = lib.DownloadJson(lib.CrawlerDirName, lib.VideoJsonPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	lib.ParseVideosFromJson(&allVideos)
	videosSortedAfterChannels = lib.SortByChannelAndShow(allVideos)
	//Connect do database
	defer lib.InitDataBaseConnection(dbConnections, "mysql", "root", "soe2020", "localhost:3306", lib.UserDBconnectionName, lib.UserDBconnectionName).Close()

	log.Print("Server has started...")
	http.Handle("/", http.FileServer(http.Dir("frontend/")))
	http.HandleFunc(lib.IncomingGetSearchRequestUrl, handleGetSearchVideos)
	http.HandleFunc(lib.IncomingGetVideosRequestUrl, handleGetAllVideos)
	http.HandleFunc(lib.IncomingGetVideosFromChannelRequestUrl, handleGetVideosByChannel)
	http.HandleFunc(lib.IncomingGetVideoClickedRequestUrl, handleGetVideoClicked)
	http.HandleFunc(lib.IncomingPostRemoveFromFavoritesRequestUrl, handleRemoveFromFavorites)
	http.HandleFunc(lib.IncomingPostUserRequestUrl, handlePostLogin)
	http.HandleFunc(lib.IncomingPostAddToFavoritesRequestUrl, handlePostAddVideoToFavorites)
	http.HandleFunc(lib.IncomingPostRegisterRequestUrl, handlePostRegisterUser)
	http.HandleFunc(lib.IncomingPostLogoutRequestUrl, handlePostLogout)
	http.HandleFunc(lib.IncomingPostCookieAUthRequestUrl, handlePostCookieAuth)
	http.HandleFunc(lib.IncomingPostCookieAUthRequestUrl, handlePostCookieAuth)
	http.HandleFunc(lib.IncomingPostFetchFavoritesRequestUrl, handlePostFetchFavorites)
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Starting Server failed: " + err.Error())
	}
}

func handlePostFetchFavorites(w http.ResponseWriter, r *http.Request) {
	log.Println("Answering handlePostFetchFavorites request started...")
	userDB := dbConnections[lib.UserDBconnectionName]
	err := userDB.Ping()
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Database connection failed: \n"+err.Error())
		return
	}
	var user lib.User
	if dErr := lib.IsUserLoggedInWithACookie(r, userDB, &user); dErr != nil {
		lib.ReportDetailedError(w, dErr)
		return
	}
	if err = lib.FillUserVideoArray(&user, userDB); err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "FillUserVideoArray failed:\n"+err.Error())
	}
	log.Println("Answered handlePostFetchFavorites request successfully")
}

func handleRemoveFromFavorites(w http.ResponseWriter, r *http.Request) {
	log.Println("Answering handleRemoveFromFavorites request started...")
	//Check connection
	userDB := dbConnections[lib.UserDBconnectionName]
	err := userDB.Ping()
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Database connection failed: \n"+err.Error())
		return
	}
	//Parse data
	err = r.ParseForm()
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Parsing request failed: \n"+err.Error())
	}
	incomingVideo := r.FormValue(lib.VideoParameter)
	if incomingVideo == "" {
		lib.ReportError(w, http.StatusBadRequest, lib.EmptyParameterResponse+lib.VideoParameter, lib.EmptyParameterResponse+lib.VideoParameter)
		return
	}
	//Check cookie
	var user lib.User
	if dErr := lib.IsUserLoggedInWithACookie(r, userDB, &user); dErr != nil {
		lib.ReportDetailedError(w, dErr)
		return
	}
	//Remove video from favorites
	result, err := userDB.Exec("delete from user_has_favorite_videos where Users_Username = ? and Video = ?", user.Username, incomingVideo)
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "sql delete failed: \n"+err.Error())
		return
	}
	returnCode, err := result.RowsAffected()
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "error in RowsAffected(): \n"+err.Error())
		return
	}
	if returnCode < 1 {
		lib.ReportError(w, http.StatusNotFound, "Video ist bei diesem User nicht favorisiert", "no rows affected by delete!")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Video aus den Favoriten entfernt!"))
	log.Println("Answered handleRemoveFromFavorites request successfully")
}

func handlePostCookieAuth(w http.ResponseWriter, r *http.Request) {
	log.Println("Answering handlePostCookieAuth request started...")
	userDB := dbConnections[lib.UserDBconnectionName]
	err := userDB.Ping()
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Database connection failed: \n"+err.Error())
		return
	}
	var user lib.User
	if dErr := lib.IsUserLoggedInWithACookie(r, userDB, &user); dErr != nil {
		lib.ReportDetailedError(w, dErr)
		return
	}
	err = lib.FillUserVideoArray(&user, userDB)
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Failed filling the favorite videos array: \n"+err.Error())
		return
	}
	userInBytes, err := json.MarshalIndent(user, "", "   ")
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Marshaling failed: \n"+err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(userInBytes)
	log.Println("Answered handlePostCookieAuth successfully")
}

func handlePostLogout(w http.ResponseWriter, r *http.Request) {
	log.Println("Answering handlePostLogout request started...")
	userDB := dbConnections[lib.UserDBconnectionName]
	err := userDB.Ping()
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Database connection failed: \n"+err.Error())
		return
	}
	var user lib.User
	if dErr := lib.IsUserLoggedInWithACookie(r, userDB, &user); dErr != nil {
		lib.ReportDetailedError(w, dErr)
		return
	}
	_, err = userDB.Exec("UPDATE users SET session_id = '0' WHERE username = ?", user.Username)
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "SQL update failed: \n"+err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Abgemeldet!"))
	log.Println("Answered handlePostLogout successfully")
}

func handleGetSearchVideos(w http.ResponseWriter, r *http.Request) {
	log.Println("Answering handleGetSearchVideos request started...")
	query := r.URL.Query()
	searchQueryResults, ok := query[lib.SearchParameter]
	if !ok || len(searchQueryResults) < 1 {
		lib.ReportError(w, http.StatusBadRequest, lib.EmptyParameterResponse+lib.SearchParameter, "Cant find parameter "+lib.VideoTitleParameter)
		return
	}
	channelQueryResults, ok := query[lib.ChannelNameParameter]
	if !ok || len(channelQueryResults) < 1 {
		lib.ReportError(w, http.StatusBadRequest, lib.EmptyParameterResponse+lib.ChannelNameParameter, "Cant find parameter "+lib.ChannelNameParameter)
		return
	}
	searchString := searchQueryResults[0]
	channelString := channelQueryResults[0]
	searchResult := make([]lib.Video, 0)
	videosFound := false
	if searchString == "" || searchString == " " {
		lib.ReportError(w, http.StatusBadRequest, lib.EmptyParameterResponse+lib.SearchParameter, lib.EmptyParameterResponse+lib.SearchParameter)
		return
	}

	if !lib.IsStringLegal(searchString) {
		lib.ReportError(w, http.StatusBadRequest, lib.IllegalParameterResponse+lib.SearchParameter, lib.IllegalParameterResponse+lib.SearchParameter)
		return
	}

	if channelString == "none" { //Search in every channel
		//Check if searchstring is a channel
		for k, v := range videosSortedAfterChannels {
			if strings.EqualFold(k, searchString) {
				for _, v2 := range v {
					for _, video := range v2 {
						searchResult = append(searchResult, video)
					}
				}
				videosFound = true
				break
			}
		}
		//Check if searchstring is a show
		if !videosFound {
			for _, v := range videosSortedAfterChannels {
				for k2, v2 := range v {
					if strings.EqualFold(k2, searchString) {
						searchResult = v2
						videosFound = true
						break
					}
				}
			}
		}
		//Search for a substring in the title
		if !videosFound {
			lowerSearchString := strings.ToLower(searchString)
			for _, v := range videosSortedAfterChannels {
				for _, v2 := range v {
					for _, video := range v2 {
						if strings.Contains(strings.ToLower(video.Title), lowerSearchString) {
							searchResult = append(searchResult, video)
							videosFound = true
						}
					}
				}
			}
		}
	} else {
		videosFromChannel := videosSortedAfterChannels[channelString]
		if videosFromChannel == nil {
			lib.ReportError(w, http.StatusNotFound, "Sender '"+channelString+"' nicht gefunden!", channelString+" doesnt exist")
			return
		}
		//Search for the show
		for k, v := range videosFromChannel {
			if strings.EqualFold(k, searchString) {
				searchResult = v
				videosFound = true
				break
			}
		}
		//Search for a substring in the title
		if !videosFound {
			lowerSearchString := strings.ToLower(searchString)
			for _, v := range videosFromChannel {
				for _, video := range v {
					if strings.Contains(strings.ToLower(video.Title), lowerSearchString) {
						searchResult = append(searchResult, video)
						videosFound = true
					}
				}
			}
		}
	}

	if !videosFound {
		log.Println("No Video found with: '" + searchString + "'!")
	}
	videosFoundJson, err := json.MarshalIndent(searchResult, "", " ")
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Marshaling failed: \n"+err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(videosFoundJson)
	log.Println("Answered handleGetSearchVideos successfully")
}

func handleGetVideoClicked(w http.ResponseWriter, r *http.Request) {
	log.Println("Answering handleGetVideoClicked request started...")
	viewCount := 1
	//Checking db connection
	userDB := dbConnections[lib.UserDBconnectionName]
	err := userDB.Ping()
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Database connection failed: \n"+err.Error())
		return
	}
	queryResults, ok := r.URL.Query()[lib.VideoTitleParameter]
	if !ok || len(queryResults) < 1 {
		lib.ReportError(w, http.StatusBadRequest, "url parameter unknown", "Cant find parameter "+lib.VideoTitleParameter)
		return
	}
	title := queryResults[0]
	log.Println("Content of parameter '" + lib.ChannelNameParameter + "': " + title)
	//Get videos from db
	rows, err := userDB.Query("select * from videos where VideoTitle = ?", title)
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Sql query failed: \n"+err.Error())
		return
	}
	if rows.Next() {
		var tempTitle string
		var tempViews int
		err = rows.Scan(&tempTitle, &tempViews)
		if err != nil {
			lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Scanning rows failed: \n"+err.Error())
			return
		}
		tempViews++
		viewCount = tempViews
		_, err = userDB.Exec("UPDATE videos SET Views = ? WHERE VideoTitle = ?", tempViews, tempTitle)
		if err != nil {
			lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "SQL update failed: \n"+err.Error())
			return
		}
	} else {
		_, err = userDB.Exec("INSERT INTO videos (VideoTitle,Views) \n Values(?,?) ", title, 1)
		if err != nil {
			viewCount = 1
			lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "SQL update failed: \n"+err.Error())
			return
		}
	}
	viewCountStr := strconv.Itoa(viewCount)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(viewCountStr))
	log.Println("Answered handleGetVideoClicked successfully")
}

func handlePostAddVideoToFavorites(w http.ResponseWriter, r *http.Request) {
	log.Println("Answering handlePostAddVideoToFavorites request started...")
	//Checking db connection
	userDB := dbConnections[lib.UserDBconnectionName]
	err := userDB.Ping()
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Database connection failed: \n"+err.Error())
		return
	}
	incomingVideo := r.FormValue("video")
	log.Println(incomingVideo)
	var user lib.User
	if dErr := lib.IsUserLoggedInWithACookie(r, userDB, &user); dErr != nil {
		lib.ReportDetailedError(w, dErr)
		return
	}
	//Check if video is already in favorites
	rows, err := userDB.Query("Select * from user_has_favorite_videos where Users_Username = ? and Video = ?", user.Username, incomingVideo)
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "SQL query checking for username and video failed: \n"+err.Error())
		return
	}
	if rows.Next() {
		lib.ReportError(w, http.StatusBadRequest, "Video ist bereits favorisiert!", "Video is already in favorites: \n")
		return
	}
	//Insert new video for the user
	_, err = userDB.Exec("INSERT INTO user_has_favorite_videos (Users_Username,Video) \n Values(?,?)", user.Username, incomingVideo)
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "SQL insert failed: \n"+err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Added video to favorites successfully"))
	log.Println("Added video to favorites successfully")
}

func handleGetVideosByChannel(w http.ResponseWriter, r *http.Request) {
	log.Println("Answering handleGetAllVideos request started...")
	queryResults, ok := r.URL.Query()[lib.ChannelNameParameter]
	if !ok || len(queryResults) < 1 {
		lib.ReportError(w, http.StatusBadRequest, lib.EmptyParameterResponse+lib.ChannelNameParameter, "Cant find parameter "+lib.ChannelNameParameter)
		return
	}
	channel := queryResults[0]
	log.Println("Content of parameter '" + lib.ChannelNameParameter + "': " + channel)
	resultSetInBytes, err := json.MarshalIndent(lib.ConvertMapToArray(videosSortedAfterChannels[channel]), "", " ")
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Marshaling failed: \n"+err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resultSetInBytes)
	log.Println("Answered handleGetAllVideos request successfully...")
}

func handlePostRegisterUser(w http.ResponseWriter, r *http.Request) {
	log.Print("answering handlePostRegisterUser request ...")

	//Checking db connection
	userDB := dbConnections[lib.UserDBconnectionName]
	err := userDB.Ping()
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Database connection failed: \n"+err.Error())
		return
	}
	//Parse username and password from request
	err = r.ParseForm()
	if err != nil {
		lib.ReportError(w, http.StatusBadRequest, "Invalid request parameters", "Parameter parsing error: "+err.Error())
		return
	}
	incomingUsername := r.FormValue(lib.UsernameParameter)
	incomingName := r.FormValue(lib.NameParameter)
	incomingPassword := r.FormValue(lib.PasswordParameter)
	//Check if any of the recieved information is empty
	if len(incomingUsername) < 1 || len(incomingName) < 1 || len(incomingPassword) < 1 {
		lib.ReportError(w, http.StatusBadRequest, lib.EmptyParameterResponse, "one or more received strings is empty\n")
		return
	}
	//Check if Username or Name is illegal
	if !lib.IsStringLegal(incomingUsername) {
		lib.ReportError(w, http.StatusBadRequest, lib.IllegalParameterResponse+lib.UsernameParameter, "forbidden chars in Username")
		return
	}
	//Get userdata from db for comparison
	rows, err := userDB.Query("select Username from users where username = ?", incomingUsername)
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Sql query failed: \n"+err.Error())
		return
	}
	//Check if username is in use
	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Scanning rows failed: \n"+err.Error())
			return
		}
		if username == incomingUsername {
			lib.ReportError(w, http.StatusBadRequest, "User already exists", "Username taken: "+username)
			return
		}
	}
	//Hash incoming password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(incomingPassword), bcrypt.MinCost)
	log.Printf("User created: Name: %s username: %s passwordhash: %s", incomingName, incomingUsername, string(passwordHash))
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Hashing password '"+incomingPassword+"' failed: \n"+err.Error())
		return
	}
	//Create user in database
	_, err = userDB.Exec("INSERT INTO users (Name,Username,PasswordHash)\nValues(?,?,?)", incomingName, incomingUsername, string(passwordHash))
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "SQL insert failed: \n"+err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Created a new User account"))
	log.Println("answered handlePostRegisterUser request successfully")
}

func handleGetAllVideos(w http.ResponseWriter, r *http.Request) {
	log.Print("Answering handleGetAllVideos request...")
	//Writing the result set to the responseWriter as a json-string
	resultSetInBytes, err := json.MarshalIndent(allVideos, "", " ")
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Marshaling failed: \n"+err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resultSetInBytes)
	log.Print("Answered handleGetAllVideos request successfully...")
}

func handlePostLogin(w http.ResponseWriter, r *http.Request) {
	log.Print("answering handlePostLogin request ...")

	//Checking db connection
	userDB := dbConnections[lib.UserDBconnectionName]
	var user lib.User
	err := userDB.Ping()
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Database connection failed: \n"+err.Error())
		return
	}

	//Parse username and password from request
	err = r.ParseForm()
	if err != nil {
		lib.ReportError(w, http.StatusBadRequest, "Invalid request parameters", "Parameter parsing error: "+err.Error())
	}
	incomingUsername := r.FormValue(lib.UsernameParameter)
	incomingPassword := r.FormValue(lib.PasswordParameter)

	if incomingUsername == "" || incomingPassword == "" {
		lib.ReportError(w, http.StatusBadRequest, lib.EmptyParameterResponse+lib.UsernameParameter+","+lib.PasswordParameter, "empty username or password")
		return
	}
	if !lib.IsStringLegal(incomingUsername) {
		lib.ReportError(w, http.StatusBadRequest, lib.IllegalParameterResponse+lib.UsernameParameter, "forbidden chars in Username")
	}

	if dErr := lib.LoginUser(userDB, &user, incomingUsername, incomingPassword); dErr != nil {
		lib.ReportDetailedError(w, dErr)
		return
	}
	if dErr := lib.PlaceCookie(w, userDB, incomingUsername); dErr != nil {
		lib.ReportDetailedError(w, dErr)
		return
	}

	err = lib.FillUserVideoArray(&user, userDB)
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Failed filling the favorite videos array: \n"+err.Error())
		return
	}
	//Writing the result set to the responseWriter as a json-string
	userInBytes, err := json.MarshalIndent(user, "", "   ")
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Marshaling failed: \n"+err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(userInBytes)
	log.Print("Answered handlePostLogin request successfully...")
}
