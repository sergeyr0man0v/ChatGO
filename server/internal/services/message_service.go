package services

import (
	"chatgo/server/internal/interfaces"
	"chatgo/server/internal/models"
	"chatgo/server/internal/util"
	"context"
	"fmt"
	"log"
	"time"
)

func (s *service) CreateMessage(c context.Context, req *interfaces.CreateMessageReq) (*interfaces.CreateMessageRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	user, err := s.Repository.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	encryptedMessage, err := util.EncryptMessage(req.Content, s.encryptKey)
	if err != nil {
		log.Printf("Failed to encrypt message: %v", err)
		return nil, fmt.Errorf("Failed to encrypt message: %v", err)
	}

	message, err := s.Repository.CreateMessage(ctx, &models.Message{
		SenderID:         user.ID,
		ChatRoomID:       req.RoomID,
		EncryptedContent: encryptedMessage,
	})
	if err != nil {
		return nil, err
	}

	return &interfaces.CreateMessageRes{
		ID:        message.ID,
		Content:   message.EncryptedContent,
		RoomID:    message.ChatRoomID,
		Username:  user.Username,
		CreatedAt: message.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *service) GetMessagesByRoomID(c context.Context, roomID string, limit int) ([]*interfaces.CreateMessageRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	messages, err := s.Repository.GetMessagesByChatRoomID(ctx, roomID, limit)
	if err != nil {
		return nil, err
	}

	result := make([]*interfaces.CreateMessageRes, len(messages))
	for i, message := range messages {
		user, err := s.Repository.GetUserByID(ctx, message.SenderID)
		if err != nil {
			return nil, err
		}
		decryptMessage, err := util.DecryptMessage(message.EncryptedContent, s.encryptKey)
		if err != nil {
			log.Printf("Failed to encrypt message: %v", err)
			return nil, fmt.Errorf("Failed to encrypt message: %v", err)
		}
		result[i] = &interfaces.CreateMessageRes{
			ID:        message.ID,
			Content:   decryptMessage,
			RoomID:    message.ChatRoomID,
			Username:  user.Username,
			CreatedAt: message.CreatedAt.Format(time.RFC3339),
		}
	}

	return result, nil
}
