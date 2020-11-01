package main

import (
	"Website/captcha"
	"Website/settings"
	"Website/views"
	"Website/ws"
	"fmt"
	"log"
	"path/filepath"

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
	router.HTMLRender = setupRenderer()

	router.GET("/", views.Index)

	authGroup := router.Group("/auth")
	{
		authGroup.GET("/login", views.LoginGET)
		authGroup.POST("/login", views.Login)

		authGroup.GET("/register", views.RegisterGET)
		authGroup.POST("/register", views.Register)

		authGroup.GET("/logout", views.Logout)
	}

	router.GET("/captcha/:id", func(c *gin.Context) {
		captcha.ShowCaptchaImage(c.Writer, c.Request, c.Param("id"))
	})

	router.Static("/static", "static")

	router.GET("/chat", func(c *gin.Context) {
		ws.ChatHandler(c.Writer, c.Request)
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
