package library

//Type of error that has enough information to answer a http request
type DetailedHttpError interface {
	Status() int         //HTTP status code
	PublicError() string //Error message that may be exposed to the client
	Error() string       //Error message that should only appear in the logs
}

//Server PublicError (statusCode >= 500)
type ServerError struct {
	statusCode         int
	publicErrorMessage string
	error              error
}

func (r *ServerError) PublicError() string {
	return r.publicErrorMessage
}
func (r *ServerError) Error() string {
	return r.error.Error()
}
func (r *ServerError) Status() int {
	return r.statusCode
}

//CLient PublicError (400 <= statusCode < 500)
type ClientError struct {
	statusCode         int
	publicErrorMessage string
	error              error
}

func (r *ClientError) PublicError() string {
	return r.publicErrorMessage
}
func (r *ClientError) Error() string {
	return r.error.Error()
}
func (r *ClientError) Status() int {
	return r.statusCode
}

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

func (r *Video) Equals(v *Video) bool {
	return r.Title == v.Title && r.Channel == v.Channel && r.Show == v.Show && r.ReleaseDate == v.ReleaseDate && r.Duration == v.Duration && r.Link == v.Link && r.PageLink == v.PageLink && v.FileName == r.FileName
}

func (r *Video) ToString() string {
	return "Channel: " + r.Channel + " Title: " + r.Title + " Show: " + r.Show + " ReleaseDate: " + r.ReleaseDate + " Duration: " + r.Duration + " Link: " + r.Link + " PageLink: " + r.PageLink + " FileName: " + r.FileName
}

type User struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Username       string `json:"username"`
	passwordHash   string
	sessionId      string
	FavoriteVideos []Video `json:"favoriteVideos"`
	ProfilePicture []byte  `json:"profilePicture"`
}

func (r *User) Equals(u *User) bool {
	if sameAttributes := r.Id == u.Id && r.Name == u.Name && r.Username == u.Username && r.passwordHash == u.passwordHash && r.sessionId == u.sessionId; !sameAttributes {
		return false
	}
	for i, v := range r.FavoriteVideos {
		if !v.Equals(&u.FavoriteVideos[i]) {
			return false
		}
	}
	return true
}

func (r *User) ToString() string {
	str := "Id: " + r.Id + " Name: " + r.Name + " Username: " + r.Username + " passwordHash: " + r.passwordHash + " sessionId: " + r.sessionId + "FavoriteVideos:\n"
	for _, v := range r.FavoriteVideos {
		str += v.ToString() + "\n"
	}
	return str
}

//Reader for tests
type DemoReader struct {
	content []byte
}

func (r *DemoReader) Read(p []byte) (n int, err error) {
	count := 0
	r.content = make([]byte, len(p))
	for _, b := range p {
		r.content[count] = b
		count++
	}
	return count, nil
}

func (r *DemoReader) GetContent() []byte {
	return r.content
}
