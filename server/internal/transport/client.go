package transport

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		c.Conn.WriteJSON(message)
	}
}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		// // Store system message about user leaving
		// _, err := hub.service.CreateMessage(context.Background(), &interfaces.CreateMessageReq{
		// 	Content:  "User left the chat",
		// 	RoomID:   c.RoomID,
		// 	Username: c.Username,
		// })
		// if err != nil {
		// 	log.Printf("Failed to store system message: %v", err)
		// }

		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var message Message
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Set RoomID from client's current room
		message.RoomID = c.RoomID

		// Validate message
		if message.Content == "" || message.Username == "" {
			log.Printf("Invalid message format: %+v", message)
			c.Conn.WriteJSON(gin.H{"error": "Invalid message format"})
			continue
		}

		// Store message in database
		// _, err = hub.service.CreateMessage(context.Background(), &interfaces.CreateMessageReq{
		// 	Content:  message.Content,
		// 	RoomID:   message.RoomID,
		// 	Username: message.Username,
		// })
		if err != nil {
			log.Printf("Failed to store message: %v", err)
			c.Conn.WriteJSON(gin.H{"error": "Failed to store message"})
			continue
		}

		// Broadcast message to room
		hub.Broadcast <- &message
	}
}
