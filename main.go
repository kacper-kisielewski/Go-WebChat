package main

import (
	"Website/settings"
	"Website/views"
	"Website/ws"
	"log"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func main() {
	router := setupRouter()
	log.Panic(router.Run(settings.Addr))
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.HTMLRender = setupRenderer()

	router.GET("/", views.Index)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", views.Login)
		authGroup.POST("/register", views.Register)
	}

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
