package services

import (
	"context"
	"server/internal/models"
	"time"
)

// CreateChatRoomReq содержит данные для создания новой чат-комнаты
type CreateChatRoomReq struct {
	ID        string              `json:"id"`
	Name      string              `json:"name"`
	Type      models.ChatRoomType `json:"type"`
	CreatorID string              `json:"creator_id"`
}

// ChatRoomRes представляет информацию о чат-комнате в ответах сервера
type ChatRoomRes struct {
	ID        string              `json:"id"`
	Name      string              `json:"name"`
	Type      models.ChatRoomType `json:"type"`
	CreatorID string              `json:"creator_id"`
	CreatedAt time.Time           `json:"created_at"`
}

// UpdateChatRoomReq содержит данные для обновления существующей чат-комнаты
type UpdateChatRoomReq struct {
	ID        string              `json:"id"`
	Name      string              `json:"name"`
	Type      models.ChatRoomType `json:"type"`
	CreatorID string              `json:"creator_id"`
}

// AddUserToChatRoomReq содержит идентификаторы пользователя и чат-комнаты для операций с участниками
type AddUserToChatRoomReq struct {
	UserID     string `json:"user_id"`
	ChatRoomID string `json:"chat_room_id"`
}

// ChatRoomService определяет методы для работы с чат-комнатами
type ChatRoomService interface {
	CreateChatRoom(c context.Context, req *CreateChatRoomReq) (*ChatRoomRes, error)
	GetChatRoomByID(c context.Context, roomID string) (*ChatRoomRes, error)
	GetChatRoomsByUserID(c context.Context, userID string) ([]*ChatRoomRes, error)
	GetAllChatRooms(c context.Context) ([]*ChatRoomRes, error)
	UpdateChatRoom(c context.Context, req *UpdateChatRoomReq) (*ChatRoomRes, error)
	DeleteChatRoom(c context.Context, roomID string) error
	AddUserToChatRoom(c context.Context, req *AddUserToChatRoomReq) error
	RemoveUserFromChatRoom(c context.Context, req *AddUserToChatRoomReq) error
}

// CreateChatRoom создает новую чат-комнату с указанными параметрами и возвращает информацию о созданной комнате
func (s *service) CreateChatRoom(c context.Context, req *CreateChatRoomReq) (*ChatRoomRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	chatRoom, err := s.Repository.CreateChatRoom(ctx, &models.ChatRoom{
		ID:        req.ID,
		Name:      req.Name,
		Type:      req.Type,
		CreatorID: req.CreatorID,
	})

	if err != nil {
		return nil, err
	}

	return &ChatRoomRes{
		ID:        chatRoom.ID,
		Name:      chatRoom.Name,
		Type:      chatRoom.Type,
		CreatorID: chatRoom.CreatorID,
		CreatedAt: chatRoom.CreatedAt,
	}, nil
}

// GetChatRoomByID возвращает информацию о чат-комнате по её идентификатору
func (s *service) GetChatRoomByID(c context.Context, roomID string) (*ChatRoomRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	chatRoom, err := s.Repository.GetChatRoomByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	return &ChatRoomRes{
		ID:        chatRoom.ID,
		Name:      chatRoom.Name,
		Type:      chatRoom.Type,
		CreatorID: chatRoom.CreatorID,
		CreatedAt: chatRoom.CreatedAt,
	}, nil
}

// GetAllChatRooms возвращает информацию о всех чат-комнатах
func (s *service) GetAllChatRooms(c context.Context) ([]*ChatRoomRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	chatRooms, err := s.Repository.GetAllChatRooms(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*ChatRoomRes, 0, len(chatRooms))
	for _, chatRoom := range chatRooms {
		result = append(result, &ChatRoomRes{
			ID:        chatRoom.ID,
			Name:      chatRoom.Name,
			Type:      chatRoom.Type,
			CreatorID: chatRoom.CreatorID,
			CreatedAt: chatRoom.CreatedAt,
		})
	}

	return result, nil
}

// GetChatRoomsByUserID возвращает список всех чат-комнат, в которых пользователь является участником
func (s *service) GetChatRoomsByUserID(c context.Context, userID string) ([]*ChatRoomRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	chatRooms, err := s.Repository.GetChatRoomsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*ChatRoomRes, 0, len(chatRooms))
	for _, chatRoom := range chatRooms {
		result = append(result, &ChatRoomRes{
			ID:        chatRoom.ID,
			Name:      chatRoom.Name,
			Type:      chatRoom.Type,
			CreatorID: chatRoom.CreatorID,
			CreatedAt: chatRoom.CreatedAt,
		})
	}

	return result, nil
}

// UpdateChatRoom обновляет существующую чат-комнату новыми данными
func (s *service) UpdateChatRoom(c context.Context, req *UpdateChatRoomReq) (*ChatRoomRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	chatRoom := &models.ChatRoom{
		ID:        req.ID,
		Name:      req.Name,
		Type:      req.Type,
		CreatorID: req.CreatorID,
	}

	updatedRoom, err := s.Repository.UpdateChatRoom(ctx, chatRoom)
	if err != nil {
		return nil, err
	}

	return &ChatRoomRes{
		ID:        updatedRoom.ID,
		Name:      updatedRoom.Name,
		Type:      updatedRoom.Type,
		CreatorID: updatedRoom.CreatorID,
		CreatedAt: updatedRoom.CreatedAt,
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
func (s *service) AddUserToChatRoom(c context.Context, req *AddUserToChatRoomReq) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	member := &models.ChatRoomMember{
		UserID:     req.UserID,
		ChatRoomID: req.ChatRoomID,
	}

	_, err := s.Repository.AddMember(ctx, member)
	return err
}

// RemoveUserFromChatRoom удаляет участника из чат-комнаты
func (s *service) RemoveUserFromChatRoom(c context.Context, req *AddUserToChatRoomReq) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	member := &models.ChatRoomMember{
		UserID:     req.UserID,
		ChatRoomID: req.ChatRoomID,
	}

	return s.Repository.DeleteMember(ctx, member)
}
