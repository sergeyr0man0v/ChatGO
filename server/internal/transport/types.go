package transport

import (
	"github.com/gorilla/websocket"
)

// Client represents a connected WebSocket client
type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string `json:"id"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

// Message represents a chat message
type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

// Room represents a chat room
type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

// RoomRes represents a chat room in responses
type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ClientRes represents a client in responses
type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
