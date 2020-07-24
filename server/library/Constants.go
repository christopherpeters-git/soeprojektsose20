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
	IncomingGetFetchProfilePictureRequestUrl  = "/getProfilePicture/"
	IncomingPostSaveProfilePictureRequestUrl  = "/setProfilePicture/"
)

//Parameter
const (
	ChannelNameParameter = "channel"
	VideoTitleParameter  = "videoTitle"
	SearchParameter      = "search"
	VideoParameter       = "video"
	UsernameParameter    = "usernameInput"
	PasswordParameter    = "passwordInput"
	NameParameter        = "nameInput"
)

//PublicError messages
const (
	InternalServerErrorResponse  = "Interner Serverfehler"
	AuthenticationFailedResponse = "Authentifizierung fehlgeschlagen"
	EmptyParameterResponse       = "Leere(r) Parameter: "
	IllegalParameterResponse     = "Illegale(r) Parameter: "
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

//miscellaneous
const (
	LetterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	AuthCookieName = "mediathekauth"
	MaxUploadSize  = 10000000 //in byte
)
