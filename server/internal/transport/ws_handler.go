package transport

import (
	"chatgo/server/internal/interfaces"
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	hub     *Hub
	service interfaces.Service
}

func NewWSHandler(h *Hub, service interfaces.Service) *WSHandler {
	return &WSHandler{
		hub:     h,
		service: service,
	}
}

// type CreateRoomReq struct {
// 	ID   string              `json:"id"`
// 	Name string              `json:"name"`
// 	Type models.ChatRoomType `json:"type"`
// }

func (h *WSHandler) CreateRoom(c *gin.Context) {
	var req interfaces.CreateChatRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from query parameters
	userID := c.Query("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	// Set user ID in context
	ctx := context.WithValue(c.Request.Context(), "user_id", userID)

	room, err := h.service.CreateChatRoom(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.hub.Rooms[room.ID] = &Room{
		ID:      room.ID,
		Name:    room.Name,
		Clients: make(map[string]*Client),
	}

	c.JSON(http.StatusOK, room)
}

// ensureDefaultRoom creates a default room if it doesn't exist
func (h *WSHandler) ensureDefaultRoom(ctx context.Context, userID string) (string, error) {
	// Check if default room exists
	rooms, err := h.service.GetAllChatRooms(ctx)
	if err != nil {
		return "", err
	}

	// Look for a room named "Default"
	for _, room := range rooms {
		if room.Name == "Default" {
			return room.ID, nil
		}
	}

	// Set user ID in context for room creation
	ctx = context.WithValue(ctx, "user_id", userID)

	// Create default room if it doesn't exist
	defaultRoom, err := h.service.CreateChatRoom(ctx, &interfaces.CreateChatRoomReq{
		Name: "Default",
	})
	if err != nil {
		return "", err
	}

	return defaultRoom.ID, nil
}

func (h *WSHandler) JoinRoom(c *gin.Context) {
	log.Printf("New WebSocket connection request")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	roomID := c.Param("roomId")
	clientID := c.Query("userId")
	username := c.Query("username")

	if clientID == "" {
		log.Printf("User ID is required")
		conn.WriteJSON(gin.H{"error": "User ID is required"})
		conn.Close()
		return
	}

	log.Printf("Client %s (username: %s) attempting to join room %s", clientID, username, roomID)

	// Handle "default" room ID
	if roomID == "default" {
		defaultRoomID, err := h.ensureDefaultRoom(c.Request.Context(), clientID)
		if err != nil {
			log.Printf("Failed to ensure default room: %v", err)
			conn.WriteJSON(gin.H{"error": "Failed to create default room"})
			conn.Close()
			return
		}
		roomID = defaultRoomID
	}

	// Convert string IDs to integers
	roomIDInt, err := strconv.ParseInt(roomID, 10, 64)
	if err != nil {
		log.Printf("Invalid room ID format: %s", roomID)
		conn.WriteJSON(gin.H{"error": "Invalid room ID format"})
		conn.Close()
		return
	}
	clientIDInt, err := strconv.ParseInt(clientID, 10, 64)
	if err != nil {
		log.Printf("Invalid user ID format: %s", clientID)
		conn.WriteJSON(gin.H{"error": "Invalid user ID format"})
		conn.Close()
		return
	}

	// Verify room exists in database
	room, err := h.service.GetChatRoomByID(c.Request.Context(), strconv.FormatInt(roomIDInt, 10))
	if err != nil {
		log.Printf("Room not found: %s", roomID)
		conn.WriteJSON(gin.H{"error": "Room not found"})
		conn.Close()
		return
	}

	// Check if user is already a member of the room
	members, err := h.service.GetMembersByChatRoomID(c.Request.Context(), strconv.FormatInt(roomIDInt, 10))
	if err != nil {
		log.Printf("Failed to get room members: %v", err)
		conn.WriteJSON(gin.H{"error": "Failed to get room members"})
		conn.Close()
		return
	}

	isMember := false
	for _, member := range members {
		if member.UserID == strconv.FormatInt(clientIDInt, 10) {
			isMember = true
			break
		}
	}

	// Only add user to room if they're not already a member
	if !isMember {
		err = h.service.AddUserToChatRoom(c.Request.Context(), &interfaces.AddUserToChatRoomReq{
			UserID:     strconv.FormatInt(clientIDInt, 10),
			ChatRoomID: strconv.FormatInt(roomIDInt, 10),
		})
		if err != nil {
			log.Printf("Failed to add user to room: %v", err)
			conn.WriteJSON(gin.H{"error": "Failed to add user to room"})
			conn.Close()
			return
		}
	}

	// Create room in hub if it doesn't exist
	if _, ok := h.hub.Rooms[roomID]; !ok {
		log.Printf("Creating new room in hub: %s", roomID)
		h.hub.Rooms[roomID] = &Room{
			ID:      roomID,
			Name:    room.Name,
			Clients: make(map[string]*Client),
		}
	}

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	log.Printf("Registering client %s in room %s", clientID, roomID)
	h.hub.Register <- cl

	m := &Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	// // Store system message in database
	// _, err = h.service.CreateMessage(c.Request.Context(), &interfaces.CreateMessageReq{
	// 	Content:  m.Content,
	// 	RoomID:   m.RoomID,
	// 	Username: m.Username,
	// })
	// if err != nil {
	// 	log.Printf("Failed to store system message: %v", err)
	// }

	log.Printf("Broadcasting join message to room %s", roomID)
	h.hub.Broadcast <- m

	go cl.writeMessage()
	cl.readMessage(h.hub)
}

func (h *WSHandler) GetAllRooms(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	rooms, err := h.service.GetAllChatRooms(c.Request.Context())
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	roomRes := make([]RoomRes, 0)
	for _, r := range rooms {
		roomRes = append(roomRes, RoomRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	conn.WriteJSON(roomRes)
}

func (h *WSHandler) GetRoomClients(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	roomId := c.Param("roomId")

	room, err := h.service.GetChatRoomByID(c.Request.Context(), roomId)
	if err != nil {
		// Return empty array instead of error object
		conn.WriteJSON(make([]ClientRes, 0))
		return
	}

	if _, ok := h.hub.Rooms[room.ID]; !ok {
		conn.WriteJSON(make([]ClientRes, 0))
		return
	}

	clients := make([]ClientRes, 0)
	for _, c := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	conn.WriteJSON(clients)
}

func (h *WSHandler) GetMessagesByRoomID(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	roomID := c.Param("roomId")
	limit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	// Handle "default" room ID
	if roomID == "default" {
		defaultRoomID, err := h.ensureDefaultRoom(c.Request.Context(), c.Query("userId"))
		if err != nil {
			conn.WriteJSON(gin.H{"error": "Failed to ensure default room"})
			return
		}
		roomID = defaultRoomID
	}

	messages, err := h.service.GetMessagesByRoomID(c.Request.Context(), roomID, limit)
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	conn.WriteJSON(messages)
}

func (h *WSHandler) GetChatRoomsByUserID(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	userID := c.Query("userId")

	res, err := h.service.GetChatRoomsByUserID(c.Request.Context(), userID)
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	conn.WriteJSON(res)
}

func (h *WSHandler) UpdateChatRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	var req interfaces.UpdateChatRoomReq
	err = conn.ReadJSON(&req)
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.UpdateChatRoom(c.Request.Context(), &req)
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	conn.WriteJSON(res)
}

func (h *WSHandler) DeleteChatRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	roomID := c.Param("roomId")

	err = h.service.DeleteChatRoom(c.Request.Context(), roomID)
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	delete(h.hub.Rooms, roomID)
	conn.WriteJSON(gin.H{"message": "Chat room deleted successfully"})
}
