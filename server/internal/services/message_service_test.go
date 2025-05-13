package services

import (
	"chatgo/server/internal/interfaces"
	"chatgo/server/internal/models"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_CreateMessage(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	req := &interfaces.CreateMessageReq{
		Content:  "Hello, world!",
		RoomID:   "room123",
		Username: "testuser",
	}

	user := &models.User{
		ID:       "user123",
		Username: req.Username,
	}

	message := &models.Message{
		ID:               "msg123",
		SenderID:         user.ID,
		ChatRoomID:       req.RoomID,
		EncryptedContent: req.Content,
		CreatedAt:        time.Now(),
	}

	mockRepo.On("GetUserByUsername", mock.Anything, req.Username).Return(user, nil)
	mockRepo.On("CreateMessage", mock.Anything, mock.MatchedBy(func(m *models.Message) bool {
		return m.SenderID == user.ID && m.ChatRoomID == req.RoomID && m.EncryptedContent == req.Content
	})).Return(message, nil)

	result, err := service.CreateMessage(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, message.ID, result.ID)
	assert.Equal(t, message.EncryptedContent, result.Content)
	assert.Equal(t, message.ChatRoomID, result.RoomID)
	assert.Equal(t, user.Username, result.Username)
	mockRepo.AssertExpectations(t)
}

func TestService_GetMessagesByRoomID(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	roomID := "room123"
	limit := 10

	user1 := &models.User{
		ID:       "user1",
		Username: "testuser1",
	}

	user2 := &models.User{
		ID:       "user2",
		Username: "testuser2",
	}

	messages := []*models.Message{
		{
			ID:               "msg1",
			SenderID:         user1.ID,
			ChatRoomID:       roomID,
			EncryptedContent: "Message 1",
			CreatedAt:        time.Now(),
		},
		{
			ID:               "msg2",
			SenderID:         user2.ID,
			ChatRoomID:       roomID,
			EncryptedContent: "Message 2",
			CreatedAt:        time.Now(),
		},
	}

	mockRepo.On("GetMessagesByChatRoomID", mock.Anything, roomID, limit).Return(messages, nil)
	mockRepo.On("GetUserByID", mock.Anything, user1.ID).Return(user1, nil)
	mockRepo.On("GetUserByID", mock.Anything, user2.ID).Return(user2, nil)

	result, err := service.GetMessagesByRoomID(context.Background(), roomID, limit)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, messages[0].ID, result[0].ID)
	assert.Equal(t, messages[0].EncryptedContent, result[0].Content)
	assert.Equal(t, user1.Username, result[0].Username)
	assert.Equal(t, messages[1].ID, result[1].ID)
	assert.Equal(t, messages[1].EncryptedContent, result[1].Content)
	assert.Equal(t, user2.Username, result[1].Username)
	mockRepo.AssertExpectations(t)
}
