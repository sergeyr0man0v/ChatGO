package models

import "time"

// ChatroomMemberRole представляет собой тип роли участника
// type Role int
type Role string

// const (
// 	Admin Role = iota
// 	Moderator
// 	Member
// )

// ChatroomMember представляет собой модель участника чата
type ChatRoomMember struct {
	UserID     int       `json:"user_id"`
	ChatRoomID int       `json:"chat_room_id"`
	JoinedAt   time.Time `json:"joined_at"`
	Role       Role      `json:"role"`
}
