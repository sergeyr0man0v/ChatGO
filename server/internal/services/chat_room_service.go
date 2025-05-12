package services

import (
	"chatgo/server/internal/interfaces"
	"chatgo/server/internal/models"
	"context"
	"fmt"
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
	CreateChatRoom(c context.Context, req *interfaces.CreateChatRoomReq) (*interfaces.CreateChatRoomRes, error)
	GetChatRoomByID(c context.Context, roomID string) (*interfaces.CreateChatRoomRes, error)
	GetChatRoomsByUserID(c context.Context, userID string) ([]*interfaces.CreateChatRoomRes, error)
	GetAllChatRooms(c context.Context) ([]*interfaces.CreateChatRoomRes, error)
	UpdateChatRoom(c context.Context, req *interfaces.UpdateChatRoomReq) (*interfaces.CreateChatRoomRes, error)
	DeleteChatRoom(c context.Context, roomID string) error
	AddUserToChatRoom(c context.Context, req *interfaces.AddUserToChatRoomReq) error
	RemoveUserFromChatRoom(c context.Context, req *interfaces.AddUserToChatRoomReq) error
	GetMembersByChatRoomID(c context.Context, roomID string) ([]*models.ChatRoomMember, error)
}

// CreateChatRoom создает новую чат-комнату с указанными параметрами и возвращает информацию о созданной комнате
func (s *service) CreateChatRoom(c context.Context, req *interfaces.CreateChatRoomReq) (*interfaces.CreateChatRoomRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// Get user from context
	userID, ok := c.Value("user_id").(string)
	if !ok {
		return nil, fmt.Errorf("user not authenticated")
	}

	chatRoom, err := s.Repository.CreateChatRoom(ctx, &models.ChatRoom{
		Name:      req.Name,
		Type:      models.Group,
		CreatorID: userID,
	})

	if err != nil {
		return nil, err
	}

	return &interfaces.CreateChatRoomRes{
		ID:   chatRoom.ID,
		Name: chatRoom.Name,
	}, nil
}

// GetChatRoomByID возвращает информацию о чат-комнате по её идентификатору
func (s *service) GetChatRoomByID(c context.Context, roomID string) (*interfaces.CreateChatRoomRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	chatRoom, err := s.Repository.GetChatRoomByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	return &interfaces.CreateChatRoomRes{
		ID:   chatRoom.ID,
		Name: chatRoom.Name,
	}, nil
}

// GetAllChatRooms возвращает информацию о всех чат-комнатах
func (s *service) GetAllChatRooms(c context.Context) ([]*interfaces.CreateChatRoomRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	chatRooms, err := s.Repository.GetAllChatRooms(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*interfaces.CreateChatRoomRes, 0, len(chatRooms))
	for _, chatRoom := range chatRooms {
		result = append(result, &interfaces.CreateChatRoomRes{
			ID:   chatRoom.ID,
			Name: chatRoom.Name,
		})
	}

	return result, nil
}

// GetChatRoomsByUserID возвращает список всех чат-комнат, в которых пользователь является участником
func (s *service) GetChatRoomsByUserID(c context.Context, userID string) ([]*interfaces.CreateChatRoomRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	chatRooms, err := s.Repository.GetChatRoomsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*interfaces.CreateChatRoomRes, 0, len(chatRooms))
	for _, chatRoom := range chatRooms {
		result = append(result, &interfaces.CreateChatRoomRes{
			ID:   chatRoom.ID,
			Name: chatRoom.Name,
		})
	}

	return result, nil
}

// UpdateChatRoom обновляет существующую чат-комнату новыми данными
func (s *service) UpdateChatRoom(c context.Context, req *interfaces.UpdateChatRoomReq) (*interfaces.CreateChatRoomRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	chatRoom := &models.ChatRoom{
		ID:   req.ID,
		Name: req.Name,
	}

	updatedRoom, err := s.Repository.UpdateChatRoom(ctx, chatRoom)
	if err != nil {
		return nil, err
	}

	return &interfaces.CreateChatRoomRes{
		ID:   updatedRoom.ID,
		Name: updatedRoom.Name,
	}, nil
}

// DeleteChatRoom удаляет чат-комнату из системы по указанному идентификатору
func (s *service) DeleteChatRoom(c context.Context, roomID string) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	chatRoom := models.ChatRoom{
		ID: roomID,
	}

	return s.Repository.DeleteChatRoom(ctx, &chatRoom)
}

// AddUserToChatRoom добавляет нового участника в существующую чат-комнату
func (s *service) AddUserToChatRoom(c context.Context, req *interfaces.AddUserToChatRoomReq) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	member := &models.ChatRoomMember{
		UserID:     req.UserID,
		ChatRoomID: req.ChatRoomID,
		MemberRole: models.Member,
		JoinedAt:   time.Now(),
	}

	_, err := s.Repository.AddMember(ctx, member)
	return err
}

// RemoveUserFromChatRoom удаляет участника из чат-комнаты
func (s *service) RemoveUserFromChatRoom(c context.Context, req *interfaces.AddUserToChatRoomReq) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	member := &models.ChatRoomMember{
		UserID:     req.UserID,
		ChatRoomID: req.ChatRoomID,
	}

	return s.Repository.DeleteMember(ctx, member)
}

// GetMembersByChatRoomID returns all members of a chat room
func (s *service) GetMembersByChatRoomID(c context.Context, roomID string) ([]*models.ChatRoomMember, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	return s.Repository.GetMembersByChatRoomID(ctx, roomID)
}
