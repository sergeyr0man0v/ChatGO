package models

import (
	"time"
)

// Message представляет собой модель сообщения
type Message struct {
	ID               string    `json:"id"`
	SenderID         string    `json:"sender_id"`
	ChatRoomID       string    `json:"chat_room_id"`
	EncryptedContent string    `json:"encrypted_content"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	IsEdited         bool      `json:"is_edited"`
}
