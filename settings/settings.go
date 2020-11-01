package settings

import (
	"os"
	"time"
)

/*
	NOTE!
	You may also edit some templates like about.html
*/

//Define global constants (settings)
const (
	Port = 8000
	Salt = 18

	SiteName         = "WebChat"
	TemplatesDir     = "templates"
	AvatarUploadsDir = "uploads/avatars"
	DefaultAvatar    = "default.jpg" //If you're changing this also make a change in models.go (db package)

	JwtTokenExpiresAt = 72000
	TokenCookieName   = "token"
	Domain            = "localhost"

	PostgresHost     = "127.0.0.1"
	PostgresUser     = "postgres"
	PostgresPassword = "password"
	PostgresDatabase = "website"
	PostgresPort     = 5432

	MaxMultipartMemory = 1 << 20

	MinimumUsernameLength         = 2
	MaximumUsernameLength         = 30
	UsernameWhitelistedCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789_"

	MaximumChatMessageLength = 140
	ChatMessageCooldown      = time.Second * 3
	ChatSystemUsername       = "@SYSTEM"

	AvatarResizeWidth  uint = 500
	AvatarResizeHeight uint = 500

	MaximumDescriptionLength = 200

	LoginInvalidCredientialsMessage = "Invalid credientials"

	ReadBufferSize  = 256
	WriteBufferSize = 256
)

//Define global variables
var (
	AvatarWhitelistedContentTypes = map[string]string{
		".jpg": "image/jpeg",
		".png": "image/png",
	}

	JwtSecret = []byte(os.Getenv("JWT_SECRET"))
)
