package interfaces

// Service defines the interface for all service operations
type Service interface {
	UserService
	MessageService
	ChatRoomService
}

// CreateUserReq represents the request to create a new user
type CreateUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateUserRes represents the response after creating a user
type CreateUserRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// LoginUserReq represents the request to login a user
type LoginUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginUserRes represents the response after logging in
type LoginUserRes struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	AccessToken string `json:"accessToken"`
}

// GetUserReq represents the request to get a user
type GetUserReq struct {
	ID string `json:"id"`
}

// GetUserRes represents the response after getting a user
type GetUserRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// CreateChatRoomReq represents the request to create a chat room
type CreateChatRoomReq struct {
	Name string `json:"name"`
}

// CreateChatRoomRes represents the response after creating a chat room
type CreateChatRoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UpdateChatRoomReq represents the request to update a chat room
type UpdateChatRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// AddUserToChatRoomReq represents the request to add a user to a chat room
type AddUserToChatRoomReq struct {
	UserID     string `json:"userId"`
	ChatRoomID string `json:"chatRoomId"`
}

// CreateMessageReq represents the request to create a message
type CreateMessageReq struct {
	Content  string `json:"content"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

// CreateMessageRes represents the response after creating a message
type CreateMessageRes struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	RoomID    string `json:"roomId"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdAt"`
}
