package websocket

import (
	"Study/websocket/pkg/jwt"
	"Study/websocket/storage"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type WebsocketHandler struct {
	storage storage.PostgresStorage
	hub *Hub
}

func NewWebsocketHandler(db *sqlx.DB) *WebsocketHandler {
	storage := *storage.New(db)
	hub, err := newHub(context.Background(), storage)
	if err != nil {
		panic(err)
	}

	go hub.run()
	return &WebsocketHandler{
		storage: storage,
		hub: hub,
	}
}
func (wbh *WebsocketHandler) ServeWebsocket(c *gin.Context) {
	token := c.Request.URL.Query().Get("Authorization")
	if token == "" {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, err := jwt.ExtractClaims(token, []byte("7VFGY5ArECvjhRU6wuLq"))
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusForbidden)
		c.Writer.Write([]byte(err.Error()))
		c.Writer.Header().Set("Content-Type", "application/json")
		return
	}
	serveWs(wbh.hub, claims["id"].(string), c.Writer, c.Request)
}
