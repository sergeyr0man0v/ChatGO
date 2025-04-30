package models

import "time"

// MemberRole представляет собой тип роли участника
type MemberRole string

const (
	Owner  MemberRole = "owner"
	Admin  MemberRole = "admin"
	Member MemberRole = "member"
)

// ChatroomMember представляет собой модель участника чата
type ChatRoomMember struct {
	UserID     string     `json:"user_id"`
	ChatRoomID string     `json:"chat_room_id"`
	JoinedAt   time.Time  `json:"joined_at"`
	MemberRole MemberRole `json:"role"`
}
