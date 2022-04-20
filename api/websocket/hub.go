package websocket

import (
	"Study/websocket/storage"
	"context"
	"fmt"

	"github.com/streadway/amqp"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients.
	messages  <-chan amqp.Delivery
	broadcast chan amqp.Delivery

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
	storage storage.PostgresStorage

}

func newHub(context context.Context, storage storage.PostgresStorage) (*Hub, error) {
	return &Hub{
		messages:     make(<-chan amqp.Delivery),
		broadcast:    make(chan amqp.Delivery),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		clients:      make(map[string]*Client),
		storage: storage,
	}, nil
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
				h.clients[client.ID] = client
			for key, val := range h.clients {
				fmt.Println(key, val, "clients")
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.send)
			}
		case message := <-h.broadcast:
			if message.Headers["session-id"] != nil {
				// if _, ok := h.clients[message.Headers["session-id"].(string)]; ok {
				// 	fmt.Println("session id", message.Headers["session-id"].(string))
				// 	select {
				// 	case h.clients[message.Headers["session-id"].(string)].send <- message.Body:
				// 		message.Ack(true)
				// 	default:
				// 		close(h.clients[message.Headers["session-id"].(string)].send)
				// 		delete(h.clients, message.Headers["session-id"].(string))
				// 	}
				// } else {
				// 	message.Ack(true)
				// }
			} else {
				message.Ack(true)
			}
		}
	}
}

