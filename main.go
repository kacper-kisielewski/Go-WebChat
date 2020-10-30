package main

import (
	"Website/settings"
	"Website/views"
	"Website/ws"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := setupRouter()
	log.Panic(router.Run(settings.Addr))
}

func setupRouter() *gin.Engine {
	router := gin.Default()
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
