package jwt

import (
	"Website/db"
	"Website/settings"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// UserClaims struct
type UserClaims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

// CreateToken returns jwt tokens
func CreateToken(username, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		username,
		email,
		jwt.StandardClaims{
			ExpiresAt: settings.JwtTokenExpiresAt,
			IssuedAt:  getCurrentUnixTime(),
		},
	})
	return token.SignedString(settings.JwtSecret)
}

// GetUsernameAndEmailFromToken extracts email from JWT token
func GetUsernameAndEmailFromToken(tokenString string) (string, string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString, new(UserClaims), func(token *jwt.Token) (interface{}, error) {
			return settings.JwtSecret, nil
		},
	)

	if claims, ok := token.Claims.(*UserClaims); ok {
		if getCurrentUnixTime() < claims.ExpiresAt {
			return "", "", errors.New("Expired token")
		}
		return claims.Username, claims.Email, nil
	}
	return "", "", err
}

// GetUserFromToken returns user model from token
func GetUserFromToken(tokenString string) (db.User, error) {
	_, email, err := GetUsernameAndEmailFromToken(tokenString)

	if err != nil {
		return db.User{}, err
	}

	return db.GetUserByEmail(email), nil
}

func getCurrentUnixTime() int64 {
	return time.Now().UTC().Unix()
}
