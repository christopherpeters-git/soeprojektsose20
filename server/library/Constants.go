package library

//URL'publicErrorMessage
const (
	IncomingPostUserRequestUrl             = "/login/"
	IncomingPostRegisterRequestUrl         = "/register/"
	IncomingPostAddToFavoritesRequestUrl   = "/addToFavorites/"
	IncomingPostLogoutRequestUrl           = "/logout/"
	IncomingPostCookieAUthRequestUrl       = "/cookieAuth/"
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

//PublicError messages
const (
	InternalServerErrorResponse  = "Internal server error - see logs"
	AuthenticationFailedResponse = "Authentication failed"
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
