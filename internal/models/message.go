package models

import (
	"database/sql"
	"time"
)

// Message представляет собой модель сообщения
type Message struct {
	ID               int           `json:"id"`
	SenderID         int           `json:"sender_id"`
	ChatRoomID       int           `json:"chat_room_id"`
	EncryptedContent string        `json:"encrypted_content"`
	ReplyToMessageID sql.NullInt64 `json:"reply_to_message_id"` // может быть NULL
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        sql.NullTime  `json:"updated_at"` // может быть NULL
	IsEdited         bool          `json:"is_edited"`
}
