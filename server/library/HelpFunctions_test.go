package library

import (
	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"
)

var exampleVideos = []Video{
	{
		Channel:     "NDR",
		Title:       "Die Nordreportage: Leben in der Jahrhundertsiedlung",
		Show:        "Die Nordreportage: Leben in der Jahrhundertsiedlung",
		ReleaseDate: "16.09.2019",
		Duration:    "00:28:31",
		Link:        "http://mediandr-a.akamaihd.net/progressive/2019/0916/TV-20190916-1113-5http.StatusBadRequest.hq.mp4",
		PageLink:    "https://www.ndr.de/fernsehen/sendungen/die_nordreportage/Leben-in-der-Jahrhundertsiedlung,sendung943144.html",
		FileName:    "76|d.mp4",
	},
	{
		Channel:     "NDR",
		Title:       "Segelschiff statt Jugendknast",
		Show:        "Letzte Chance an Bord",
		ReleaseDate: "02.11.2019",
		Duration:    "00:29:36",
		Link:        "http://mediandr-a.akamaihd.net/progressive/2018/0412/TV-20180412-1528-4http.StatusBadRequest.hq.mp4",
		PageLink:    "https://www.ndr.de/fernsehen/sendungen/die_reportage/Segelschiff-statt-Jugendknast,sendung610984.html",
		FileName:    "76|d.mp4",
	}}

var userColumns = []string{"Id", "Name", "Username", "PasswordHash", "Session_Id"}

func TestIsStringLegal(t *testing.T) {
	legalString := "Hallo123"
	illegalStrings := [4]string{"Hal<lo", "Hallo>", "hal/lo", "hall.o"}
	//Tests, if legal string is legal
	if !IsStringLegal(legalString) {
		t.Errorf("string '%s' should be legal!", legalString)
	}
	//Tests, if illegal strings are legal
	for _, str := range illegalStrings {
		if IsStringLegal(str) {
			t.Errorf("string '%s' should be illegal!", legalString)
		}
	}
}

func TestFillUserVideoArray(t *testing.T) {
	db, mock, err := sqlmock.New()
	columns := []string{"Users_Username", "Video"}
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	givenUser := User{
		Id:             "1",
		Name:           "Bob",
		Username:       "bob123",
		FavoriteVideos: nil,
	}

	expectedUser := User{
		Id:             "1",
		Name:           "Bob",
		Username:       "bob123",
		FavoriteVideos: exampleVideos,
	}
	favVideo1 := "{\n  \"channel\": \"NDR\",\n  \"title\": \"Die Nordreportage: Leben in der Jahrhundertsiedlung\",\n  \"show\": \"Die Nordreportage: Leben in der Jahrhundertsiedlung\",\n  \"releaseDate\": \"16.09.2019\",\n  \"duration\": \"00:28:31\",\n  \"link\": \"http://mediandr-a.akamaihd.net/progressive/2019/0916/TV-20190916-1113-5http.StatusBadRequest.hq.mp4\",\n  \"pageLink\": \"https://www.ndr.de/fernsehen/sendungen/die_nordreportage/Leben-in-der-Jahrhundertsiedlung,sendung943144.html\",\n  \"fileName\": \"76|d.mp4\"\n }"
	favVideo2 := "{\n  \"channel\": \"NDR\",\n  \"title\": \"Segelschiff statt Jugendknast\",\n  \"show\": \"Letzte Chance an Bord\",\n  \"releaseDate\": \"02.11.2019\",\n  \"duration\": \"00:29:36\",\n  \"link\": \"http://mediandr-a.akamaihd.net/progressive/2018/0412/TV-20180412-1528-4http.StatusBadRequest.hq.mp4\",\n  \"pageLink\": \"https://www.ndr.de/fernsehen/sendungen/die_reportage/Segelschiff-statt-Jugendknast,sendung610984.html\",\n  \"fileName\": \"76|d.mp4\"\n }"
	resultRows := sqlmock.NewRows(columns).AddRow(givenUser.Username, favVideo1).AddRow(givenUser.Username, favVideo2)
	mock.ExpectQuery("select (.+) from user_has_favorite_videos where Users_Username = (.+)").WillReturnRows(resultRows)
	//Tests, if FillUservideoArray throws no error as expected
	if err := FillUserVideoArray(&givenUser, db); err != nil {
		t.Error("Unexpected error \n" + err.Error())
		return
	} else if !givenUser.Equals(&expectedUser) { //Tests, if user is as expected after the FillUservideoArray call
		t.Error("Users are not equal!: givenUser: " + givenUser.ToString() + " expectedUser: " + expectedUser.ToString())
	}
}

func TestConvertMapToArray(t *testing.T) {
	exampleMap := make(map[string]map[string][]Video)
	exampleMap["abc"] = make(map[string][]Video)
	exampleMap["abc"]["def"] = make([]Video, 0)
	exampleMap["abc"]["def"] = append(exampleMap["abc"]["def"], exampleVideos[0])
	exampleMap["abc"]["ghi"] = make([]Video, 0)
	exampleMap["abc"]["ghi"] = append(exampleMap["abc"]["ghi"], exampleVideos[1])

	result := ConvertMapToArray(exampleMap["abc"])
	//Tests if an the result of ConvertMapToArray is as expected
	for i, v := range result {
		if !v.Equals(&exampleVideos[i]) {
			t.Errorf("Video didnt match: resultVideo: %s expectedVideo: %s\n", v.ToString(), exampleVideos[i].ToString())
		}
	}
}

func TestSortByChannelAndShow(t *testing.T) {
	exampleMap := make(map[string]map[string][]Video)
	exampleMap[exampleVideos[0].Channel] = make(map[string][]Video)
	exampleMap[exampleVideos[0].Channel][exampleVideos[0].Show] = make([]Video, 0)
	exampleMap[exampleVideos[0].Channel][exampleVideos[0].Show] = append(exampleMap[exampleVideos[0].Channel][exampleVideos[0].Show], exampleVideos[0])
	exampleMap[exampleVideos[0].Channel][exampleVideos[1].Show] = make([]Video, 0)
	exampleMap[exampleVideos[0].Channel][exampleVideos[1].Show] = append(exampleMap[exampleVideos[0].Channel][exampleVideos[1].Show], exampleVideos[1])
	//Tests, if SortByChannelAndShow creates a sorted map as expected
	if !reflect.DeepEqual(SortByChannelAndShow(exampleVideos), exampleMap) {
		t.Error("Maps are not the same")
	}
}

func TestLoginUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	//Test if user is able to be logged in
	var givenUser User
	incomingUsername := "maxMustermann"
	incomingPassword := "muster123"
	incomingPasswordHash, err := bcrypt.GenerateFromPassword([]byte(incomingPassword), bcrypt.MinCost)
	givenSessionId := "0"
	if err != nil {
		t.Fatalf("error '%s' was not expected creating a hash", err.Error())
	}
	mock.ExpectQuery("select (.+) from users where username = ?").WillReturnRows(sqlmock.NewRows(userColumns).AddRow("0", "Max Mustermann", incomingUsername, string(incomingPasswordHash), givenSessionId))
	expectedUser := User{
		Id:             "0",
		Name:           "Max Mustermann",
		Username:       incomingUsername,
		passwordHash:   string(incomingPasswordHash),
		sessionId:      "0",
		FavoriteVideos: nil,
	}
	//Tests if LoginUser thorws no error as expected
	if dErr := LoginUser(db, &givenUser, incomingUsername, incomingPassword); dErr != nil {
		t.Error("login failed unexpected: " + dErr.Error())
	} else if !givenUser.Equals(&expectedUser) { //Tests if logged in user is the expecte user
		t.Errorf("given user didnt match expected user: givenUser: %s expectedUser: %s\n", givenUser.ToString(), expectedUser.ToString())
	}
	mock.ExpectQuery("select (.+) from users where username = ?").WillReturnRows(sqlmock.NewRows(userColumns).AddRow("0", "Max Mustermann", incomingUsername, string(incomingPasswordHash), givenSessionId))
	var givenUser2 User
	wrongPassword := "wrongPassword"
	dErr := LoginUser(db, &givenUser2, incomingUsername, wrongPassword)
	//Test if user cant log with wrong password in as expected
	if dErr == nil {
		t.Errorf("Error expected!")
	} else if dErr.Status() != http.StatusForbidden {
		t.Errorf("Expected error status http.StatusForbidden, got: " + strconv.FormatInt(int64(dErr.Status()), 10) + " " + dErr.Error())
	}
}

func TestIsUserLoggedInWithACookie(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var givenUser User
	//Creating a demo request and cookie
	exampleSessionId := GenerateSessionId(255)
	expire := time.Now().AddDate(0, 0, 2)
	sessionCookie := http.Cookie{
		Name:       AuthCookieName,
		Value:      exampleSessionId,
		Path:       "/",
		Domain:     "localhost",
		Expires:    expire,
		RawExpires: expire.Format(time.UnixDate),
		MaxAge:     172800,
		Secure:     false,
		HttpOnly:   true,
		SameSite:   http.SameSiteLaxMode,
	}
	reader := DemoReader{}
	request, err := http.NewRequest("POST", "localhost", &reader)
	if err != nil {
		t.Errorf("unexpected error while creating a request: \n %s", err.Error())
	}

	expectedUser := User{
		Id:             "1",
		Name:           "Bob",
		Username:       "bob123",
		sessionId:      exampleSessionId,
		passwordHash:   "abcdef",
		FavoriteVideos: nil,
	}
	//Test if login throws error if no cookie is there
	if dErr := IsUserLoggedInWithACookie(request, db, &givenUser); dErr == nil {
		t.Errorf("login was successful unexpectedly!")
	}
	//Test if login is successful with right session id
	request, err = http.NewRequest("POST", "localhost", &reader)
	if err != nil {
		t.Errorf("unexpected error while creating a request: \n %s", err.Error())
		return
	}
	//Test if login is successful with right session id
	request.AddCookie(&sessionCookie)
	mock.ExpectQuery("select (.+) from users where session_id = ?").WillReturnRows(sqlmock.NewRows(userColumns).AddRow(expectedUser.Id, expectedUser.Name, expectedUser.Username, expectedUser.passwordHash, expectedUser.sessionId))
	if dErr := IsUserLoggedInWithACookie(request, db, &givenUser); dErr != nil {
		t.Errorf("login threw unexepcted error: \n %s", dErr.Error())
	} else if !givenUser.Equals(&expectedUser) {
		t.Errorf("given user does not match expected user: gU:\n %s\n eU:\n %s\n", givenUser, expectedUser)
	}
	//Test if login throws error with a wrong session id
	mock.ExpectQuery("select (.+) from users where session_id = ?").WillReturnRows(sqlmock.NewRows(userColumns))
	sessionCookie.Value = "wrongSessionId"
	if dErr := IsUserLoggedInWithACookie(request, db, &givenUser); dErr == nil {
		t.Errorf("login was successful unexpectedly!")
	}
}
