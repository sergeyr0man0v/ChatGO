package services

import (
	"context"
	"server/internal/models"
	"time"
)

type CreateMessageReq struct {
	Content  string `json:"content"`
	RoomID   string `json:"room_id"`
	Username string `json:"username"`
}

type MessageRes struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	RoomID    string    `json:"room_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageService interface {
	CreateMessage(c context.Context, req *CreateMessageReq) (*MessageRes, error)
	GetMessagesByRoomID(c context.Context, roomID string, limit int) ([]*MessageRes, error)
}

func (s *service) CreateMessage(c context.Context, req *CreateMessageReq) (*MessageRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	user, err := s.Repository.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	message, err := s.Repository.CreateMessage(ctx, &models.Message{
		SenderID:         user.ID,
		ChatRoomID:       req.RoomID,
		EncryptedContent: req.Content,
	})
	if err != nil {
		return nil, err
	}

	return &MessageRes{
		ID:        message.ID,
		Content:   message.EncryptedContent,
		RoomID:    message.ChatRoomID,
		Username:  user.Username,
		CreatedAt: message.CreatedAt,
	}, nil
}

func (s *service) GetMessagesByRoomID(c context.Context, roomID string, limit int) ([]*MessageRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	messages, err := s.Repository.GetMessagesByChatRoomID(ctx, roomID, limit)
	if err != nil {
		return nil, err
	}

	result := make([]*MessageRes, len(messages))
	for i, message := range messages {
		user, err := s.Repository.GetUserByID(ctx, message.SenderID)
		if err != nil {
			return nil, err
		}
		result[i] = &MessageRes{
			ID:        message.ID,
			Content:   message.EncryptedContent,
			RoomID:    message.ChatRoomID,
			Username:  user.Username,
			CreatedAt: message.CreatedAt,
		}
	}

	return result, nil
}
