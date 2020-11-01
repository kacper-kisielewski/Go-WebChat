package views

import (
	"Website/db"
	"Website/jwt"
	"Website/settings"
	"log"
	"net/http"

	"Website/captcha"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//Index view
func Index(c *gin.Context) {
	renderTemplate(c, "index", map[string]interface{}{
		"tokenCookieName":    settings.TokenCookieName,
		"chatSystemUsername": settings.ChatSystemUsername,
		"domain":             settings.Domain,
		"port":               settings.Port,
	})
}

//LoginGET view
func LoginGET(c *gin.Context) {
	renderTemplate(c, "login", nil)
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

	c.SetCookie(
		settings.TokenCookieName,
		token,
		settings.JwtTokenExpiresAt,
		"/",
		settings.Domain,
		true,
		false,
	)
	c.Redirect(http.StatusFound, "/")
}

//RegisterGET view
func RegisterGET(c *gin.Context) {
	renderTemplate(c, "register", map[string]interface{}{
		"captchaID": captcha.NewCaptcha(),
	})
}

//Register view
func Register(c *gin.Context) {
	var registerBody RegisterBody

	if !captcha.VerifyCaptcha(c.PostForm("captcha_id"), c.PostForm("captcha_solution")) {
		c.String(http.StatusForbidden, "Invalid captcha")
		return
	}

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

	c.Redirect(http.StatusFound, "/auth/login")
}

//Logout view
func Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     settings.TokenCookieName,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	c.Redirect(http.StatusMovedPermanently, "/")
	return
}
