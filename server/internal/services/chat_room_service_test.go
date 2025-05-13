package services

import (
	"chatgo/server/internal/interfaces"
	"chatgo/server/internal/models"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_CreateChatRoom(t *testing.T) {
	testCases := []struct {
		name        string
		req         *interfaces.CreateChatRoomReq
		ctx         context.Context
		mockSetup   func(mockRepo *MockRepository)
		expectError bool
		checkResult func(t *testing.T, result *interfaces.CreateChatRoomRes, err error)
	}{
		{
			name: "Successfully create chat room",
			req: &interfaces.CreateChatRoomReq{
				Name: "Test Room",
			},
			ctx: context.WithValue(context.Background(), "user_id", "user123"),
			mockSetup: func(mockRepo *MockRepository) {
				expectedChatRoom := &models.ChatRoom{
					ID:        "room123",
					Name:      "Test Room",
					Type:      models.Group,
					CreatorID: "user123",
					CreatedAt: time.Now(),
				}

				mockRepo.On("CreateChatRoom", mock.Anything, mock.MatchedBy(func(chatRoom *models.ChatRoom) bool {
					return chatRoom.Name == "Test Room" && chatRoom.Type == models.Group && chatRoom.CreatorID == "user123"
				})).Return(expectedChatRoom, nil)
			},
			expectError: false,
			checkResult: func(t *testing.T, result *interfaces.CreateChatRoomRes, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "room123", result.ID)
				assert.Equal(t, "Test Room", result.Name)
			},
		},
		{
			name: "Missing user ID in context",
			req: &interfaces.CreateChatRoomReq{
				Name: "Test Room",
			},
			ctx:         context.Background(),
			mockSetup:   func(mockRepo *MockRepository) {},
			expectError: true,
			checkResult: func(t *testing.T, result *interfaces.CreateChatRoomRes, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
			},
		},
		{
			name: "Repository error",
			req: &interfaces.CreateChatRoomReq{
				Name: "Test Room",
			},
			ctx: context.WithValue(context.Background(), "user_id", "user123"),
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("CreateChatRoom", mock.Anything, mock.MatchedBy(func(chatRoom *models.ChatRoom) bool {
					return chatRoom.Name == "Test Room" && chatRoom.Type == models.Group && chatRoom.CreatorID == "user123"
				})).Return(nil, errors.New("database error"))
			},
			expectError: true,
			checkResult: func(t *testing.T, result *interfaces.CreateChatRoomRes, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			service := NewService(mockRepo)

			tc.mockSetup(mockRepo)

			result, err := service.CreateChatRoom(tc.ctx, tc.req)

			tc.checkResult(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetChatRoomByID(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	roomID := "room123"
	expectedChatRoom := &models.ChatRoom{
		ID:        roomID,
		Name:      "Test Room",
		Type:      models.Group,
		CreatorID: "user123",
		CreatedAt: time.Now(),
	}

	mockRepo.On("GetChatRoomByID", mock.Anything, roomID).Return(expectedChatRoom, nil)

	result, err := service.GetChatRoomByID(context.Background(), roomID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedChatRoom.ID, result.ID)
	assert.Equal(t, expectedChatRoom.Name, result.Name)
	mockRepo.AssertExpectations(t)
}

func TestService_GetAllChatRooms(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	expectedChatRooms := []*models.ChatRoom{
		{
			ID:        "room1",
			Name:      "Room 1",
			Type:      models.Group,
			CreatorID: "user1",
			CreatedAt: time.Now(),
		},
		{
			ID:        "room2",
			Name:      "Room 2",
			Type:      models.Direct,
			CreatorID: "user2",
			CreatedAt: time.Now(),
		},
	}

	mockRepo.On("GetAllChatRooms", mock.Anything).Return(expectedChatRooms, nil)

	result, err := service.GetAllChatRooms(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, len(expectedChatRooms))
	assert.Equal(t, expectedChatRooms[0].ID, result[0].ID)
	assert.Equal(t, expectedChatRooms[1].ID, result[1].ID)
	mockRepo.AssertExpectations(t)
}

func TestService_GetChatRoomsByUserID(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	userID := "user123"
	expectedChatRooms := []*models.ChatRoom{
		{
			ID:        "room1",
			Name:      "Room 1",
			Type:      models.Group,
			CreatorID: userID,
			CreatedAt: time.Now(),
		},
		{
			ID:        "room2",
			Name:      "Room 2",
			Type:      models.Direct,
			CreatorID: "user2",
			CreatedAt: time.Now(),
		},
	}

	mockRepo.On("GetChatRoomsByUserID", mock.Anything, userID).Return(expectedChatRooms, nil)

	result, err := service.GetChatRoomsByUserID(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, len(expectedChatRooms))
	assert.Equal(t, expectedChatRooms[0].ID, result[0].ID)
	assert.Equal(t, expectedChatRooms[1].ID, result[1].ID)
	mockRepo.AssertExpectations(t)
}

func TestService_UpdateChatRoom(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	req := &interfaces.UpdateChatRoomReq{
		ID:   "room123",
		Name: "Updated Room",
	}

	expectedChatRoom := &models.ChatRoom{
		ID:        req.ID,
		Name:      req.Name,
		Type:      models.Group,
		CreatorID: "user123",
		CreatedAt: time.Now(),
	}

	mockRepo.On("UpdateChatRoom", mock.Anything, mock.MatchedBy(func(chatRoom *models.ChatRoom) bool {
		return chatRoom.ID == req.ID && chatRoom.Name == req.Name
	})).Return(expectedChatRoom, nil)

	result, err := service.UpdateChatRoom(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedChatRoom.ID, result.ID)
	assert.Equal(t, expectedChatRoom.Name, result.Name)
	mockRepo.AssertExpectations(t)
}

func TestService_DeleteChatRoom(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	roomID := "room123"

	mockRepo.On("DeleteChatRoom", mock.Anything, mock.MatchedBy(func(chatRoom *models.ChatRoom) bool {
		return chatRoom.ID == roomID
	})).Return(nil)

	err := service.DeleteChatRoom(context.Background(), roomID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_AddUserToChatRoom(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	req := &interfaces.AddUserToChatRoomReq{
		UserID:     "user123",
		ChatRoomID: "room123",
	}

	expectedMember := &models.ChatRoomMember{
		UserID:     req.UserID,
		ChatRoomID: req.ChatRoomID,
		JoinedAt:   time.Now(),
		MemberRole: models.Member,
	}

	mockRepo.On("AddMember", mock.Anything, mock.MatchedBy(func(member *models.ChatRoomMember) bool {
		return member.UserID == req.UserID && member.ChatRoomID == req.ChatRoomID
	})).Return(expectedMember, nil)

	err := service.AddUserToChatRoom(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_RemoveUserFromChatRoom(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	req := &interfaces.AddUserToChatRoomReq{
		UserID:     "user123",
		ChatRoomID: "room123",
	}

	mockRepo.On("DeleteMember", mock.Anything, mock.MatchedBy(func(member *models.ChatRoomMember) bool {
		return member.UserID == req.UserID && member.ChatRoomID == req.ChatRoomID
	})).Return(nil)

	err := service.RemoveUserFromChatRoom(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_GetMembersByChatRoomID(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	roomID := "room123"
	expectedMembers := []*models.ChatRoomMember{
		{
			UserID:     "user1",
			ChatRoomID: roomID,
			JoinedAt:   time.Now(),
			MemberRole: models.Owner,
		},
		{
			UserID:     "user2",
			ChatRoomID: roomID,
			JoinedAt:   time.Now(),
			MemberRole: models.Member,
		},
	}

	mockRepo.On("GetMembersByChatRoomID", mock.Anything, roomID).Return(expectedMembers, nil)

	result, err := service.GetMembersByChatRoomID(context.Background(), roomID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, len(expectedMembers))
	assert.Equal(t, expectedMembers[0].UserID, result[0].UserID)
	assert.Equal(t, expectedMembers[1].UserID, result[1].UserID)
	mockRepo.AssertExpectations(t)
}
