package settings

import "os"

//Define global constants (settings)
const (
	Addr = ":8000"
	Salt = 18

	SiteName     = "WebChat"
	TemplatesDir = "templates"

	MaximumUsernameLength         = 30
	MinimumUsernameLength         = 2
	UsernameWhitelistedCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789_"

	RegisterSuccessfullMessage = "User registered"

	LoginInvalidCredientialsMessage = "Invalid credientials"

	ReadBufferSize  = 1024
	WriteBufferSize = 1024
)

//Define global variables
var (
	JwtSecret = []byte(os.Getenv("JWT_SECRET"))
)
