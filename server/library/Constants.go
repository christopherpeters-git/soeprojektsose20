package library

//URL'publicErrorMessage
const (
	IncomingPostUserRequestUrl                = "/login/"
	IncomingPostRegisterRequestUrl            = "/register/"
	IncomingPostAddToFavoritesRequestUrl      = "/addToFavorites/"
	IncomingPostRemoveFromFavoritesRequestUrl = "/removeFromFavorites/"
	IncomingPostLogoutRequestUrl              = "/logout/"
	IncomingPostCookieAUthRequestUrl          = "/cookieAuth/"
	IncomingGetSearchRequestUrl               = "/search"
	IncomingGetVideosFromChannelRequestUrl    = "/getVideoByChannel"
	IncomingGetVideoClickedRequestUrl         = "/clickVideo"
	IncomingGetVideosRequestUrl               = "/getVideos/"
)

//Parameter
const (
	ChannelNameParameter = "channel"
	VideoTitleParameter  = "videoTitle"
	SearchParameter      = "search"
	VideoParameter       = "video"
	UsernameParameter    = "usernameInput"
	PasswordParameter    = "passwordInput"
)

//PublicError messages
const (
	InternalServerErrorResponse  = "Interner Serverfehler"
	AuthenticationFailedResponse = "Authentifizierung fehlgeschlagen"
	EmptyParameterResponse       = "Leere(r) Parameter:"
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
