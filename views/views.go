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
	"github.com/lestrrat-go/strftime"
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
	renderTemplate(c, "login", nil, "Login")
}

//Login view
func Login(c *gin.Context) {
	var loginBody LoginBody

	if err := c.ShouldBindWith(&loginBody, binding.Form); err != nil {
		c.String(http.StatusUnprocessableEntity, err.Error())
		return
	}

	user := db.GetUserByEmail(loginBody.Email)
	if !db.AuthenticateUser(user.HashedPassword, []byte(loginBody.Password)) {
		c.String(http.StatusForbidden, settings.LoginInvalidCredientialsMessage)
		return
	}

	token, err := jwt.CreateToken(user.Username, user.Email)
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
	}, "Register")
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
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     settings.TokenCookieName,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	c.Redirect(http.StatusMovedPermanently, "/")
}

//Profile view
func Profile(c *gin.Context) {
	username := c.Param("username")
	user := db.GetUserByUsername(username)
	createdAt, _ := strftime.Format("%x", user.CreatedAt)
	userExists := (user.Username != "")

	renderTemplate(c, "profile", map[string]interface{}{
		"username":   user.Username,
		"createdAt":  createdAt,
		"userExists": userExists,
	}, (func() string {
		if userExists {
			return username + "'s Profile"
		}
		return "User not found"
	})())
}
