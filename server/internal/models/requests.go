package models

// AddUserToChatRoomReq содержит идентификаторы пользователя и чат-комнаты для операций с участниками
type AddUserToChatRoomReq struct {
	UserID     string `json:"user_id"`
	ChatRoomID string `json:"chat_room_id"`
}

// CreateMessageReq содержит данные для создания нового сообщения
type CreateMessageReq struct {
	Content  string `json:"content"`
	RoomID   string `json:"room_id"`
	Username string `json:"username"`
}

// CreateChatRoomReq содержит данные для создания новой чат-комнаты
type CreateChatRoomReq struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Type      ChatRoomType `json:"type"`
	CreatorID string       `json:"creator_id"`
}

// UpdateChatRoomReq содержит данные для обновления существующей чат-комнаты
type UpdateChatRoomReq struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Type      ChatRoomType `json:"type"`
	CreatorID string       `json:"creator_id"`
}
