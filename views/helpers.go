package views

import (
	"Website/jwt"
	"Website/settings"
	"net/http"

	"github.com/gin-gonic/gin"
)

//AuthenticateContext func
func AuthenticateContext(c *gin.Context) error {
	tokenCookie, err := c.Cookie(settings.TokenCookieName)
	if err != nil {
		return err
	}

	_, err = jwt.GetEmailFromToken(tokenCookie)
	return err
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
	data["isAuthenticated"] = (AuthenticateContext(c) == nil)

	for key, value := range obj {
		data[key] = value
	}
	c.HTML(http.StatusOK, name+".html", data)
}
