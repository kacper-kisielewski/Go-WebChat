package settings

import (
	"os"
	"time"
)

//Define global constants (settings)
const (
	Addr = ":8000"
	Salt = 18

	SiteName     = "WebChat"
	TemplatesDir = "templates"

	JwtTokenExpiresAt = 72000
	TokenCookieName   = "token"
	Domain            = "localhost"

	MaximumUsernameLength         = 30
	MinimumUsernameLength         = 2
	UsernameWhitelistedCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789_"

	MaximumChatMessageLength = 100
	ChatMessageCooldown      = time.Second * 3
	ChatSystemUsername       = "@SYSTEM"

	LoginInvalidCredientialsMessage = "Invalid credientials"

	ReadBufferSize  = 128
	WriteBufferSize = 128
)

//Define global variables
var (
	JwtSecret = []byte(os.Getenv("JWT_SECRET"))
)
