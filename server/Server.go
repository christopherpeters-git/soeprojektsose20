package main

import (
	lib "./library"
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"image/png"
	"io/ioutil"
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
	http.HandleFunc(lib.IncomingGetFetchProfilePictureRequestUrl, handleGetFetchProfilePicture)
	http.HandleFunc(lib.IncomingGetLogoutRequestUrl, handleGetLogout)
	http.HandleFunc(lib.IncomingGetCookieAuthRequestUrl, handleGetCookieAuth)
	http.HandleFunc(lib.IncomingGetFetchFavoritesRequestUrl, handleGetFetchFavorites)
	http.HandleFunc(lib.IncomingPostRemoveFromFavoritesRequestUrl, handlePostRemoveFromFavorites)
	http.HandleFunc(lib.IncomingPostUserRequestUrl, handlePostLogin)
	http.HandleFunc(lib.IncomingPostAddToFavoritesRequestUrl, handlePostAddVideoToFavorites)
	http.HandleFunc(lib.IncomingPostRegisterRequestUrl, handlePostRegisterUser)
	http.HandleFunc(lib.IncomingPostSaveProfilePictureRequestUrl, handlePostSaveProfilePicture)
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Starting Server failed: " + err.Error())
	}
}

//Sends all favorite videos in a json-array for the authenticated user
func handleGetFetchFavorites(w http.ResponseWriter, r *http.Request) {
	log.Println("Answering handleGetFetchFavorites request started...")
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
		return
	}
	videosInBytes, err := json.MarshalIndent(user.FavoriteVideos, "", " ")
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "marshaling failed: \n"+err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(videosInBytes)
	log.Println("Answered handleGetFetchFavorites request successfully")
}

//saves the incoming picture in the db for the authenticated user
func handlePostSaveProfilePicture(w http.ResponseWriter, r *http.Request) {
	log.Println("Answering handlePostSaveProfilePicture request started...")
	userDB := dbConnections[lib.UserDBconnectionName]
	//Checking DB connection
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
	if err = r.ParseMultipartForm(lib.MaxUploadSize); err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Parsing request form failed: \n"+err.Error())
		return
	}
	//Parsing FormData File
	_, header, err := r.FormFile("profilepicture")
	if err != nil {
		lib.ReportError(w, http.StatusBadRequest, "Keinen Wert für 'profilepicture' gefunden", "error getting value for profilepicture: \n"+err.Error())
		return
	}
	if header.Size == 0 || header.Size > lib.MaxUploadSize {
		lib.ReportError(w, http.StatusBadRequest, "Bild ist zu groß!\n Die Grenze beträgt "+strconv.FormatInt(lib.MaxUploadSize, 10)+"bytes", "file size too large: "+strconv.FormatInt(header.Size, 10))
		return
	}
	//Open file and trying to decode it to png
	imageFile, err := header.Open()
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "opening file failed: \n"+err.Error())
		return
	}
	_, err = png.Decode(imageFile)
	if err != nil {
		lib.ReportError(w, http.StatusBadRequest, "Datei ist kein PNG-Bild!", "Decoding png failed: \n"+err.Error())
		return
	}
	//Open and save file to DB
	imageFile, err = header.Open()
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "opening file failed: \n"+err.Error())
		return
	}
	user.ProfilePicture, err = ioutil.ReadAll(imageFile)
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "reading file failed: \n"+err.Error())
		return
	}
	log.Printf("SAVED PROFILE PIC: %d \n", len(user.ProfilePicture))
	_, err = userDB.Exec("UPDATE users SET profile_picture = ? WHERE username = ?", user.ProfilePicture, user.Username)
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "SQL update failed: \n"+err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Profile picture set"))
	log.Println("Answered handlePostSaveProfilePicture request successfully")
}

//Sends the profile picture for the authenticated user or the standard profile picture
func handleGetFetchProfilePicture(w http.ResponseWriter, r *http.Request) {
	log.Println("Answering handleGetFetchProfilePicture request started...")
	userDB := dbConnections[lib.UserDBconnectionName]
	err := userDB.Ping()
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Database connection failed: \n"+err.Error())
		return
	}
	var user lib.User
	if dErr := lib.IsUserLoggedInWithACookie(r, userDB, &user); dErr != nil {
		log.Println("authentication failed, loading standard picture")
		user.ProfilePicture, err = ioutil.ReadFile(lib.StandardAvatarPath)
		if err != nil {
			lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Reading avatar.png failed: \n"+err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(user.ProfilePicture)
		return
	}
	rows, err := userDB.Query("select profile_picture from users where username = ?", user.Username)
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "sql query failed: \n"+err.Error())
		return
	}
	if rows.Next() {
		err := rows.Scan(&user.ProfilePicture)
		if err != nil {
			lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "scanning failed: \n"+err.Error())
			return
		}
		if user.ProfilePicture == nil {
			log.Println("no picture found, loading standard picture")
			user.ProfilePicture, err = ioutil.ReadFile(lib.StandardAvatarPath)
			if err != nil {
				lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Reading avatar.png failed: \n"+err.Error())
				return
			}
		}
	} else {
		log.Println("no user found, loading standard picture")
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Logged in user not found!")
		return
	}
	log.Printf("picture byte array len: %d\n", len(user.ProfilePicture))
	w.WriteHeader(http.StatusOK)
	w.Write(user.ProfilePicture)
	log.Println("Answered handleGetFetchProfilePicture request successfully")
}

//Removes the incoming video for the authenticated user
func handlePostRemoveFromFavorites(w http.ResponseWriter, r *http.Request) {
	log.Println("Answering handlePostRemoveFromFavorites request started...")
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
	log.Println("Answered handlePostRemoveFromFavorites request successfully")
}

//Sends userinformation as a json-object for the authenticated user
func handleGetCookieAuth(w http.ResponseWriter, r *http.Request) {
	log.Println("Answering handleGetCookieAuth request started...")
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
	userInBytes, err := json.MarshalIndent(user, "", "   ")
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Marshaling failed: \n"+err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(userInBytes)
	log.Println("Answered handleGetCookieAuth successfully")
}

//Sets SessionID to 0 for the authenticated user
func handleGetLogout(w http.ResponseWriter, r *http.Request) {
	log.Println("Answering handleGetLogout request started...")
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
	log.Println("Answered handleGetLogout successfully")
}

//Returns found videos for a specified channel in as a json
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

//Increase the viewcount for an incoming videoTitle or creates a new entry for the title. Returns the new view count
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

//Adds an incoming video to the favorites of the authenticated user
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
	rows, err := userDB.Query("Select users_username,video from user_has_favorite_videos where users_username = ? and video = ?", user.Username, incomingVideo)
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
	w.Write([]byte("Video hinzugefügt!"))
	log.Println("Added video to favorites successfully")
}

//Sends all videos in an array for an incoming channel
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

//creates a new user in the db for the send credentials
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
	//Check if Username is illegal
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
	for rows.Next() {
		lib.ReportError(w, http.StatusBadRequest, "User already exists", "Username taken: "+incomingUsername)
		return
	}
	//Hash incoming password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(incomingPassword), bcrypt.MinCost)
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "Hashing password '"+incomingPassword+"' failed: \n"+err.Error())
		return
	}
	log.Printf("User created: Name: %s username: %s passwordhash: %s", incomingName, incomingUsername, string(passwordHash))
	//Create user in database
	_, err = userDB.Exec("INSERT INTO users (Name,Username,PasswordHash)\nValues(?,?,?)", incomingName, incomingUsername, string(passwordHash))
	if err != nil {
		lib.ReportError(w, http.StatusInternalServerError, lib.InternalServerErrorResponse, "SQL insert failed: \n"+err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Neuer Account angelegt"))
	log.Println("answered handlePostRegisterUser request successfully")
}

//Sends all videos as a json
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

//Sends userinformation for an authenticated user, places cookie in client if successful
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
