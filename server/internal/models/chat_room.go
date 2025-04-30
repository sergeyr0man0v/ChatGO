package models

import "time"

// ChatRoomType представляет собой тип чата
type ChatRoomType string

const (
	Direct ChatRoomType = "direct"
	Group  ChatRoomType = "group"
)

// ChatRoom представляет собой модель чата
type ChatRoom struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Type      ChatRoomType `json:"type"`
	CreatedAt time.Time    `json:"created_at"`
	CreatorID string       `json:"creator_id"`
}
