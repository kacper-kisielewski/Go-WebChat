package settings

import (
	"os"
	"regexp"
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

	PostgresHost     = "db"
	PostgresUser     = "postgres"
	PostgresPassword = "password"
	PostgresDatabase = "website"
	PostgresPort     = 5432

	MaxMultipartMemory = 1 << 20

	MainChannel = "general"

	MinimumUsernameLength         = 2
	MaximumUsernameLength         = 30
	UsernameWhitelistedCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789_"

	MinimumEmailLength = 4

	MaximumChannelNameLength = 10
	MinimumChannelNameLength = 3

	MaximumChatMessageLength = 140
	ChatMessageCooldown      = time.Second * 3
	ChatSystemUsername       = "@SYSTEM"

	AvatarJPGQuality        = 40
	AvatarResizeWidth  uint = 500
	AvatarResizeHeight uint = 500

	MaximumDescriptionLength = 200

	LoginInvalidCredientialsMessage = "Invalid credientials"
	UserDisabledMessage             = "User disabled"

	ReadBufferSize  = 256
	WriteBufferSize = 256
)

//Define global variables
var (
	PinnedChannels = []string{
		"OffTopic", "Coding",
	}

	ChannelNameRegex = regexp.MustCompile(`[^\w]`)

	AvatarWhitelistedContentTypes = map[string]string{
		".jpg": "image/jpeg",
		".png": "image/png",
	}

	JwtSecret = []byte(os.Getenv("JWT_SECRET"))
)
