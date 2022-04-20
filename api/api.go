package api

import (
	"net/http"

	"Study/websocket/api/handler"
	webSkt "Study/websocket/api/websocket"
	"Study/websocket/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)


func New(db *sqlx.DB, cfg config.Config) *gin.Engine {
	r := gin.New()

	r.Static("/images", "./static/images")

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowHeaders = append(config.AllowHeaders, "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	config.AllowHeaders = append(config.AllowHeaders, "*")

	handler := handler.New(db, cfg)
	websockerHandler := webSkt.NewWebsocketHandler(db)

	r.Use(cors.New(config))
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "Api gateway"})
	})
	{
		r.GET("/ws", websockerHandler.ServeWebsocket)
		r.POST("/login", handler.Login)
		r.GET("/chat-users/:id", handler.GetChatUsers)
		r.GET("/chat-history", handler.GetChatHistory)
		r.POST("/message", handler.CreateMessage)
		r.POST("/update-read", handler.UpdateReadStatus)
	}

	return r
}