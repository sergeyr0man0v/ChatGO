package services

import (
	"server/internal/models"
	"time"
)

// // Request/response DTOs for chat room member operations
// type AddMemberReq struct {
// 	UserID     string `json:"user_id"`
// 	ChatRoomID string `json:"chat_room_id"`
// 	Role       string `json:"role"`
// }

// type ChatRoomMemberRes struct {
// 	UserID     string    `json:"user_id"`
// 	Username   string    `json:"username,omitempty"`
// 	ChatRoomID string    `json:"chat_room_id"`
// 	Role       string    `json:"role"`
// 	JoinedAt   time.Time `json:"joined_at"`
// }

// type UpdateMemberRoleReq struct {
// 	ID   string `json:"id"`
// 	Role string `json:"role"`
// }

// // // Chat room member service interface
// // type ChatRoomMemberService interface {
// // 	AddMember(c context.Context, req *AddMemberReq) (*ChatRoomMemberRes, error)
// // 	GetMembersByRoomID(c context.Context, roomID string) ([]*ChatRoomMemberRes, error)
// // 	GetMemberByID(c context.Context, memberID string) (*ChatRoomMemberRes, error)
// // 	UpdateMemberRole(c context.Context, req *UpdateMemberRoleReq) (*ChatRoomMemberRes, error)
// // 	RemoveMember(c context.Context, memberID string) error
// // }

// Common service interface that combines all specific service interfaces
type Service interface {
	UserService
	MessageService
	ChatRoomService
	// ChatRoomMemberService
}

type service struct {
	models.Repository
	timeout time.Duration
}

func NewService(repository models.Repository) Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
	}
}
