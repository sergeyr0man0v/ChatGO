package services

// import (
// 	"chatgo/server/internal/interfaces"
// 	"chatgo/server/internal/models"
// 	"chatgo/server/internal/util"
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// func TestService_CreateMessage(t *testing.T) {
// 	mockRepo := new(MockRepository)
// 	service := NewService(mockRepo, config)

// 	req := &interfaces.CreateMessageReq{
// 		Content:  "Hello, world!",
// 		RoomID:   "room123",
// 		Username: "testuser",
// 	}

// 	user := &models.User{
// 		ID:       "user123",
// 		Username: req.Username,
// 	}

// 	encryptedMessage, err := util.EncryptMessage(req.Content, config.encryptKey)
// 	assert.NoError(t, err)

// 	message := &models.Message{
// 		ID:               "msg123",
// 		SenderID:         user.ID,
// 		ChatRoomID:       req.RoomID,
// 		EncryptedContent: encryptedMessage,
// 		CreatedAt:        time.Now(),
// 	}

// 	mockRepo.On("GetUserByUsername", mock.Anything, req.Username).Return(user, nil)
// 	mockRepo.On("CreateMessage", mock.Anything, mock.MatchedBy(func(m *models.Message) bool {
// 		return m.SenderID == user.ID && m.ChatRoomID == req.RoomID && m.EncryptedContent == encryptedMessage
// 	})).Return(message, nil)

// 	result, err := service.CreateMessage(context.Background(), req)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, result)
// 	assert.Equal(t, message.ID, result.ID)
// 	assert.Equal(t, encryptedMessage, result.Content)
// 	assert.Equal(t, message.ChatRoomID, result.RoomID)
// 	assert.Equal(t, user.Username, result.Username)
// 	mockRepo.AssertExpectations(t)
// }

// func TestService_GetMessagesByRoomID(t *testing.T) {
// 	mockRepo := new(MockRepository)
// 	service := NewService(mockRepo, config)

// 	roomID := "room123"
// 	limit := 10

// 	user1 := &models.User{
// 		ID:       "user1",
// 		Username: "testuser1",
// 	}

// 	user2 := &models.User{
// 		ID:       "user2",
// 		Username: "testuser2",
// 	}

// 	encryptedMsg1, err := util.EncryptMessage("Message 1", config.encryptKey)
// 	assert.NoError(t, err)
// 	encryptedMsg2, err := util.EncryptMessage("Message 2", config.encryptKey)
// 	assert.NoError(t, err)

// 	messages := []*models.Message{
// 		{
// 			ID:               "msg1",
// 			SenderID:         user1.ID,
// 			ChatRoomID:       roomID,
// 			EncryptedContent: encryptedMsg1,
// 			CreatedAt:        time.Now(),
// 		},
// 		{
// 			ID:               "msg2",
// 			SenderID:         user2.ID,
// 			ChatRoomID:       roomID,
// 			EncryptedContent: encryptedMsg2,
// 			CreatedAt:        time.Now(),
// 		},
// 	}

// 	mockRepo.On("GetMessagesByChatRoomID", mock.Anything, roomID, limit).Return(messages, nil)
// 	mockRepo.On("GetUserByID", mock.Anything, user1.ID).Return(user1, nil)
// 	mockRepo.On("GetUserByID", mock.Anything, user2.ID).Return(user2, nil)

// 	result, err := service.GetMessagesByRoomID(context.Background(), roomID, limit)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, result)
// 	assert.Len(t, result, 2)

// 	// First message assertions
// 	assert.Equal(t, messages[0].ID, result[0].ID)
// 	decryptedMsg1, err := util.DecryptMessage(messages[0].EncryptedContent, config.encryptKey)
// 	assert.NoError(t, err)
// 	assert.Equal(t, decryptedMsg1, result[0].Content)
// 	assert.Equal(t, user1.Username, result[0].Username)

// 	// Second message assertions
// 	assert.Equal(t, messages[1].ID, result[1].ID)
// 	decryptedMsg2, err := util.DecryptMessage(messages[1].EncryptedContent, config.encryptKey)
// 	assert.NoError(t, err)
// 	assert.Equal(t, decryptedMsg2, result[1].Content)
// 	assert.Equal(t, user2.Username, result[1].Username)

// 	mockRepo.AssertExpectations(t)
// }
