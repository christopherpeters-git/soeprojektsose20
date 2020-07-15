package library

//URL's
const (
	IncomingPostUserRequestUrl             = "/login/"
	IncomingPostRegisterRequestUrl         = "/register/"
	IncomingPostAddToFavoritesRequestUrl   = "/addToFavorites/"
	IncomingPostLogoutRequestUrl           = "/logout/"
	IncomingGetSearchRequestUrl            = "/search"
	IncomingGetVideosFromChannelRequestUrl = "/getVideoByChannel"
	IncomingGetVideoClickedRequestUrl      = "/clickVideo"
	IncomingGetVideosRequestUrl            = "/getVideos/"
)

//Parameter
const (
	ChannelNameParameter = "channel"
	VideoTitleParameter  = "videoTitle"
	VideoSearchParameter = "search"
)

//Error messages
const (
	InternalServerErrorResponse = "Internal server error - see logs"
)

//db connection names
const (
	UserDBconnectionName = "userdb"
)

//Paths
const (
	CrawlerDirName = "crawler"
	VideoJsonPath  = CrawlerDirName + "/good.json"
)

//Characters allowed in SessionID
const (
	LetterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	AuthCookieName = "mediathekauth"
)
