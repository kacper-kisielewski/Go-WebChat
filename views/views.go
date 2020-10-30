package views

import (
	"Website/db"
	"Website/jwt"
	"Website/settings"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Index view
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Index")
}

//Login view
func Login(c *gin.Context) {
	var (
		email    = c.PostForm("email")
		password = []byte(c.PostForm("password"))
	)
	if !db.AuthenticateUser(email, password) {
		c.String(http.StatusForbidden, settings.LoginInvalidCredientialsMessage)
		return
	}

	token, err := jwt.CreateToken(email)
	if err != nil {
		log.Panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
	})
}

//Register view
func Register(c *gin.Context) {
	var (
		username = c.PostForm("username")
		email    = c.PostForm("email")
		password = c.PostForm("password")
	)

	if err := db.RegisterUser(username, email, password); err != nil {
		log.Panic(err)
	}

	c.String(http.StatusOK, settings.RegisterSuccessfullMessage)
}
