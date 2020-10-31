package views

import (
	"Website/db"
	"Website/jwt"
	"Website/settings"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//Index view
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Index")
}

//Login view
func Login(c *gin.Context) {
	var loginBody LoginBody

	if err := c.ShouldBindWith(&loginBody, binding.Form); err != nil {
		c.String(http.StatusUnprocessableEntity, err.Error())
		return
	}

	if !db.AuthenticateUser(loginBody.Email, []byte(loginBody.Password)) {
		c.String(http.StatusForbidden, settings.LoginInvalidCredientialsMessage)
		return
	}

	token, err := jwt.CreateToken(loginBody.Email)
	if err != nil {
		log.Panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
	})
}

//Register view
func Register(c *gin.Context) {
	var registerBody RegisterBody

	if err := c.ShouldBindWith(&registerBody, binding.Form); err != nil {
		c.String(http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err := db.RegisterUser(
		registerBody.Username,
		registerBody.Email,
		[]byte(registerBody.Password),
	); err != nil {
		c.String(http.StatusUnprocessableEntity, err.Error())
		return
	}

	c.String(http.StatusOK, settings.RegisterSuccessfullMessage)
}
