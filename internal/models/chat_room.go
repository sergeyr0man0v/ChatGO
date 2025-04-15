package models

import "time"

// ChatRoomType представляет собой тип чата
type ChatRoomType int

const (
	Direct ChatRoomType = iota
	Group
)

// ChatRoom представляет собой модель чата
type ChatRoom struct {
	ID        int          `json:"id"`
	Name      string       `json:"name"`
	Type      ChatRoomType `json:"type"`
	CreatedAt time.Time    `json:"created_at"`
	CreatorID int          `json:"creator_id"`
}
