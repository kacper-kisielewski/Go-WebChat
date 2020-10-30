package settings

import "os"

//Define global constants (settings)
const (
	Addr = ":8000"
	Salt = 18

	MaximumPasswordLength = 30
	MinimumPasswordLength = 2

	RegisterSuccessfullMessage = "User registered"

	LoginInvalidCredientialsMessage = "Invalid credientials"

	ReadBufferSize  = 1024
	WriteBufferSize = 1024
)

//Define global variables
var (
	JwtSecret = []byte(os.Getenv("JWT_SECRET"))
)
