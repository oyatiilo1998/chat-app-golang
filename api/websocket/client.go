package websocket

import (
	"Study/websocket/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	// "anor/anor_websocket/modules/anor/response"
	// "anor/anor_websocket/modules/anor/tracking_service"
	// "anor/anor_websocket/pkg/logger"
	// "anor/anor_websocket/pkg/pubsub"

	"github.com/gorilla/websocket"
)

//

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 500 * time.Millisecond

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 250 * time.Millisecond

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	ID string
	hub       *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var req models.Message
		err := c.conn.ReadJSON(&req)
		fmt.Println(req.To, c.hub.clients)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println(err)
			}
			fmt.Println(err, "err")
			break
		}

		resp, err := c.hub.storage.CreateMessage(req)
		if err != nil {
			fmt.Println(err)
		}

		message, err := json.Marshal(resp)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(resp)

		c.hub.clients[req.To].send <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			var resp models.Message
			json.Unmarshal(message, &resp)
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				fmt.Println(ok, "ok")
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteJSON(resp)
			if err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, ID string, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		ID: ID,
		hub:       hub,
		conn:      conn,
		send:      make(chan []byte, 256),
	}
	
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
