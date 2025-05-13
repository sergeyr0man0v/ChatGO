package services

import (
	"chatgo/server/internal/interfaces"
	"chatgo/server/internal/models"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_CreateUser(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	req := &interfaces.CreateUserReq{
		Username: "testuser",
		Password: "password123",
	}

	createdUser := &models.User{
		ID:                "user123",
		Username:          req.Username,
		EncryptedPassword: "hashedpassword",
	}

	mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
		return u.Username == req.Username && u.EncryptedPassword != req.Password
	})).Return(createdUser, nil)

	result, err := service.CreateUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdUser.ID, result.ID)
	assert.Equal(t, createdUser.Username, result.Username)
	mockRepo.AssertExpectations(t)
}

func TestService_Login(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	req := &interfaces.LoginUserReq{
		Username: "testuser",
		Password: "password123",
	}

	hashedPassword := "$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8iw0hLyhsiG"

	user := &models.User{
		ID:                "user123",
		Username:          req.Username,
		EncryptedPassword: hashedPassword,
	}

	mockRepo.On("GetUserByUsername", mock.Anything, req.Username).Return(user, nil)

	result, err := service.Login(context.Background(), req)

	if err == nil {
		assert.NotNil(t, result)
		assert.Equal(t, user.ID, result.ID)
		assert.Equal(t, user.Username, result.Username)
		assert.NotEmpty(t, result.AccessToken)
	}
	mockRepo.AssertExpectations(t)
}

func TestService_GetUserByID(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	userID := "user123"
	req := &interfaces.GetUserReq{
		ID: userID,
	}

	user := &models.User{
		ID:       userID,
		Username: "testuser",
	}

	mockRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)

	result, err := service.GetUserByID(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Username, result.Username)
	mockRepo.AssertExpectations(t)
}

func TestService_GetAllUsers(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	users := []*models.User{
		{
			ID:       "user1",
			Username: "testuser1",
		},
		{
			ID:       "user2",
			Username: "testuser2",
		},
	}

	mockRepo.On("GetAllUsers", mock.Anything).Return(users, nil)

	result, err := service.GetAllUsers(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, len(users))

	for i, user := range users {
		assert.Equal(t, user.ID, result[i].ID)
		assert.Equal(t, user.Username, result[i].Username)
	}

	mockRepo.AssertExpectations(t)
}
