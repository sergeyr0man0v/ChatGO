package models

import (
	"context"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetAllUsers(ctx context.Context) ([]*User, error)
}

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *Message) (*Message, error)
	GetMessageByID(ctx context.Context, messageID string) (*Message, error)
	GetMessagesByChatRoomID(ctx context.Context, roomID string, limit int) ([]*Message, error)
}

type ChatRoomRepository interface {
	CreateChatRoom(ctx context.Context, chatRoom *ChatRoom) (*ChatRoom, error)
	GetChatRoomByID(ctx context.Context, chatRoomID string) (*ChatRoom, error)
	GetChatRoomsByUserID(ctx context.Context, userID string) ([]*ChatRoom, error)
	GetMembersByChatRoomID(ctx context.Context, chatRoomID string) ([]*ChatRoomMember, error)
	GetAllChatRooms(ctx context.Context) ([]*ChatRoom, error)
	UpdateChatRoom(ctx context.Context, chatRoom *ChatRoom) (*ChatRoom, error)
	UpdateMemberRole(ctx context.Context, member *ChatRoomMember) (*ChatRoomMember, error)
	DeleteChatRoom(ctx context.Context, chatRoom *ChatRoom) error
	AddMember(ctx context.Context, member *ChatRoomMember) (*ChatRoomMember, error)
	DeleteMember(ctx context.Context, member *ChatRoomMember) error
	// DeleteMembersByChatRoomID(ctx context.Context, /**sql.Tx*/ chatRoomID string) error
}

// type ChatRoomMemberRepository interface {
// 	AddMember(ctx context.Context, member *ChatRoomMember) (*ChatRoomMember, error)
// 	GetMembersByUserID(ctx context.Context, userID string) ([]*ChatRoomMember, error)
// 	GetMembersByChatRoomID(ctx context.Context, chatRoomID string) ([]*ChatRoomMember, error)
// 	GetMemberByUserAndRoomID(ctx context.Context, userID string, chatRoomID string) (*ChatRoomMember, error)
// 	UpdateMemberRole(ctx context.Context, member *ChatRoomMember) (*ChatRoomMember, error)
// 	RemoveMember(ctx context.Context, member *ChatRoomMember) error
// 	RemoveMembersByUserID(ctx context.Context, userID string) error
// 	RemoveMembersByChatRoomID(ctx context.Context, chatRoomID string) error
// }

type Repository interface {
	UserRepository
	MessageRepository
	ChatRoomRepository
	//ChatRoomMemberRepository
}
