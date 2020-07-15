package library

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
	Id             string `json:"id"`
	Name           string `json:"name"`
	Username       string `json:"username"`
	passwordHash   string
	sessionId      string
	FavoriteVideos []Video `json:"favoriteVideos"`
}
