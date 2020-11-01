package views

import (
	"Website/jwt"
	"Website/settings"
	"net/http"

	"github.com/gin-gonic/gin"
)

//AuthenticateContext func
func AuthenticateContext(c *gin.Context) (string, string, error) {
	tokenCookie, err := c.Cookie(settings.TokenCookieName)
	if err != nil {
		return "", "", err
	}

	return jwt.GetUsernameAndEmailFromToken(tokenCookie)
}

func renderTemplate(c *gin.Context, name string, obj map[string]interface{}, title ...string) {
	var data gin.H

	if len(title) > 0 {
		data = gin.H{
			"title": title[0] + " | " + settings.SiteName,
		}
	} else {
		data = gin.H{
			"title": settings.SiteName,
		}
	}
	data["siteName"] = settings.SiteName

	username, _, err := AuthenticateContext(c)
	data["isAuthenticated"] = (err == nil)
	data["currentUsername"] = username

	for key, value := range obj {
		data[key] = value
	}
	c.HTML(http.StatusOK, name+".html", data)
}
