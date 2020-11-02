package views

import (
	"Website/db"
	"Website/jwt"
	"Website/settings"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"Website/captcha"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/lestrrat-go/strftime"
)

//Index view
func Index(c *gin.Context) {
	if IsAuthenticated(c) {
		c.Redirect(http.StatusPermanentRedirect, "/channel/"+settings.MainChannel)
		return
	}

	renderTemplate(c, "index", map[string]interface{}{
		"mainChannel": settings.MainChannel,
	})
}

//Channel view
func Channel(c *gin.Context) {
	if !IsAuthenticated(c) {
		c.Status(http.StatusForbidden)
		return
	}

	var (
		channel     = strings.ToLower(c.Param("channel"))
		isNameValid = IsValidChannelName(channel)
	)

	renderTemplate(c, "channel", map[string]interface{}{
		"tokenCookieName":    settings.TokenCookieName,
		"chatSystemUsername": settings.ChatSystemUsername,
		"channel":            channel,
		"domain":             settings.Domain,
		"port":               settings.Port,
		"isNameValid":        isNameValid,
	}, func() string {
		if isNameValid {
			return "#" + channel
		}
		return "Invalid channel"
	}())
}

//About view
func About(c *gin.Context) {
	renderTemplate(c, "about", nil, "About")
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
	var (
		username     = c.Param("username")
		user         = db.GetUserByUsername(username)
		createdAt, _ = strftime.Format("%x", user.CreatedAt)
		userExists   = (user.Username != "")
	)

	renderTemplate(c, "profile", map[string]interface{}{
		"username":    user.Username,
		"createdAt":   createdAt,
		"avatar":      user.Avatar,
		"description": user.Description,
		"userExists":  userExists,
		"isOwner": (func() bool {
			if currentUsername, _, _ := AuthenticateContext(c); currentUsername == user.Username {
				return true
			}
			return false
		})(),
	}, (func() string {
		if userExists {
			return username + "'s Profile"
		}
		return "User not found"
	})())
}

//EditDescriptionGET view
func EditDescriptionGET(c *gin.Context) {
	if !IsAuthenticated(c) {
		c.Status(http.StatusForbidden)
		return
	}

	user := GetUserFromContext(c)
	renderTemplate(c, "edit_desc", map[string]interface{}{
		"description": user.Description,
	}, "Edit Description")
}

//EditDescription view
func EditDescription(c *gin.Context) {
	user := GetUserFromContext(c)
	description := c.PostForm("description")

	if len(description) > settings.MaximumDescriptionLength {
		c.String(http.StatusUnprocessableEntity, fmt.Sprintf(
			"Description cannot exceed %d characters", settings.MaximumDescriptionLength,
		))
		return
	}

	db.DB.Model(&user).Update("description", strings.TrimSpace(description))
	c.Redirect(http.StatusFound, "/profile/"+user.Username)
}

//EditAvatarGET view
func EditAvatarGET(c *gin.Context) {
	if !IsAuthenticated(c) {
		c.Status(http.StatusForbidden)
		return
	}

	renderTemplate(c, "edit_avatar", nil, "Edit Avatar")
}

//EditAvatar view
func EditAvatar(c *gin.Context) {
	var avatarDir = func(fileName string) string {
		return path.Join(settings.AvatarUploadsDir, fileName)
	}

	user := GetUserFromContext(c)
	avatar, err := c.FormFile("avatar")
	if err != nil {
		log.Panic(err)
	}

	ext, ok := getContentTypeExt(avatar.Header.Get("Content-Type"))
	if !ok {
		c.String(http.StatusUnprocessableEntity, "Invalid content type")
		return
	}

	fileName := generateAvatarFileName(ext)
	err = resizeAndSave(avatar, avatarDir(fileName))
	if err != nil {
		log.Panic(err)
	}

	if user.Avatar != settings.DefaultAvatar {
		if err = os.Remove(avatarDir(user.Avatar)); err != nil {
			log.Println(err)
		}
	}

	db.DB.Model(&user).Update("avatar", fileName)

	c.Redirect(http.StatusFound, "/profile/"+user.Username)
}
