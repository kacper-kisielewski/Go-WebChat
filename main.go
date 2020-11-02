package main

import (
	"Website/captcha"
	"Website/settings"
	"Website/views"
	"Website/ws"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func init() {
	log.Println("Starting HTTP server on port", settings.Port)
}

func main() {
	router := setupRouter()
	log.Panic(router.Run(fmt.Sprintf(":%d", settings.Port)))
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = settings.MaxMultipartMemory
	router.HTMLRender = setupRenderer()

	router.GET("/", views.Index)
	router.GET("/about", views.About)

	authGroup := router.Group("/auth")
	{
		authGroup.GET("/login", views.LoginGET)
		authGroup.POST("/login", views.Login)

		authGroup.GET("/register", views.RegisterGET)
		authGroup.POST("/register", views.Register)

		authGroup.GET("/logout", views.Logout)
	}

	router.GET("/profile/:username", views.Profile)

	settingsGroup := router.Group("/settings")
	{
		settingsGroup.GET("/desc", views.EditDescriptionGET)
		settingsGroup.POST("/desc", views.EditDescription)

		settingsGroup.GET("/avatar", views.EditAvatarGET)
		settingsGroup.POST("/avatar", views.EditAvatar)
	}

	router.GET("/captcha/:id", func(c *gin.Context) {
		captcha.ShowCaptchaImage(c.Writer, c.Request, c.Param("id"))
	})

	router.GET("/channel/:channel", views.Channel)

	router.Static("/avatars", settings.AvatarUploadsDir)
	router.Static("/static", "static")

	router.GET("/chat/:channel", func(c *gin.Context) {
		ws.ChatHandler(c.Writer, c.Request, strings.ToLower(c.Param("channel")))
	})

	return router
}

func setupRenderer() multitemplate.Renderer {
	renderer := multitemplate.NewRenderer()

	layouts, err := templatesGlob("layouts")
	if err != nil {
		log.Panic(err)
	}

	includes, err := templatesGlob("includes")

	for _, include := range includes {
		files := append(layouts, include)
		renderer.AddFromFiles(filepath.Base(include), files...)
	}

	return renderer
}

func templatesGlob(folder string) ([]string, error) {
	return filepath.Glob(
		filepath.Join(
			settings.TemplatesDir, folder, "*.html",
		),
	)
}
