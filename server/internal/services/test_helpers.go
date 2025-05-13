package services

import (
	"chatgo/server/internal/models"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of models.Repository
type MockRepository struct {
	mock.Mock
}

// Implement necessary methods from models.Repository interface
func (m *MockRepository) CreateChatRoom(ctx context.Context, chatRoom *models.ChatRoom) (*models.ChatRoom, error) {
	args := m.Called(ctx, chatRoom)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatRoom), args.Error(1)
}
func (m *MockRepository) GetChatRoomByID(ctx context.Context, chatRoomID string) (*models.ChatRoom, error) {
	args := m.Called(ctx, chatRoomID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatRoom), args.Error(1)
}
func (m *MockRepository) GetAllChatRooms(ctx context.Context) ([]*models.ChatRoom, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.ChatRoom), args.Error(1)
}
func (m *MockRepository) GetChatRoomsByUserID(ctx context.Context, userID string) ([]*models.ChatRoom, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.ChatRoom), args.Error(1)
}
func (m *MockRepository) UpdateChatRoom(ctx context.Context, chatRoom *models.ChatRoom) (*models.ChatRoom, error) {
	args := m.Called(ctx, chatRoom)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatRoom), args.Error(1)
}
func (m *MockRepository) DeleteChatRoom(ctx context.Context, chatRoom *models.ChatRoom) error {
	args := m.Called(ctx, chatRoom)
	return args.Error(0)
}

func (m *MockRepository) AddMember(ctx context.Context, member *models.ChatRoomMember) (*models.ChatRoomMember, error) {
	args := m.Called(ctx, member)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatRoomMember), args.Error(1)
}

func (m *MockRepository) DeleteMember(ctx context.Context, member *models.ChatRoomMember) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MockRepository) GetMembersByChatRoomID(ctx context.Context, chatRoomID string) ([]*models.ChatRoomMember, error) {
	args := m.Called(ctx, chatRoomID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.ChatRoomMember), args.Error(1)
}

// Additional mock methods for user service tests
func (m *MockRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

// Additional mock methods for message service tests
func (m *MockRepository) CreateMessage(ctx context.Context, message *models.Message) (*models.Message, error) {
	args := m.Called(ctx, message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Message), args.Error(1)
}

func (m *MockRepository) GetMessageByID(ctx context.Context, messageID string) (*models.Message, error) {
	args := m.Called(ctx, messageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Message), args.Error(1)
}

func (m *MockRepository) GetMessagesByChatRoomID(ctx context.Context, roomID string, limit int) ([]*models.Message, error) {
	args := m.Called(ctx, roomID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Message), args.Error(1)
}

func (m *MockRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) UpdateMemberRole(ctx context.Context, member *models.ChatRoomMember) (*models.ChatRoomMember, error) {
	args := m.Called(ctx, member)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatRoomMember), args.Error(1)
}
