package interfaces

import (
	"chatgo/server/internal/models"
	"context"
	"time"
)

// ChatRoomRes представляет информацию о чат-комнате в ответах сервера
type ChatRoomRes struct {
	ID        string              `json:"id"`
	Name      string              `json:"name"`
	Type      models.ChatRoomType `json:"type"`
	CreatorID string              `json:"creator_id"`
	CreatedAt time.Time           `json:"created_at"`
}

// ChatRoomService определяет методы для работы с чат-комнатами
type ChatRoomService interface {
	CreateChatRoom(c context.Context, req *CreateChatRoomReq) (*CreateChatRoomRes, error)
	GetChatRoomByID(c context.Context, roomID string) (*CreateChatRoomRes, error)
	GetChatRoomsByUserID(c context.Context, userID string) ([]*CreateChatRoomRes, error)
	GetAllChatRooms(c context.Context) ([]*CreateChatRoomRes, error)
	UpdateChatRoom(c context.Context, req *UpdateChatRoomReq) (*CreateChatRoomRes, error)
	DeleteChatRoom(c context.Context, roomID string) error
	AddUserToChatRoom(c context.Context, req *AddUserToChatRoomReq) error
	RemoveUserFromChatRoom(c context.Context, req *AddUserToChatRoomReq) error
	GetMembersByChatRoomID(c context.Context, roomID string) ([]*models.ChatRoomMember, error)
}
