package transport

import (
	"net/http"
	"server/internal/models"
	"server/internal/services"
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
	service services.Service
}

func NewWSHandler(h *Hub, service *services.Service) *WSHandler {
	return &WSHandler{
		hub:     h,
		service: *service,
	}
}

type CreateRoomReq struct {
	ID   string              `json:"id"`
	Name string              `json:"name"`
	Type models.ChatRoomType `json:"type"`
}

func (h *WSHandler) CreateRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	var req CreateRoomReq
	err = conn.ReadJSON(&req)
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	room, err := h.service.CreateChatRoom(c.Request.Context(), &services.CreateChatRoomReq{
		ID:        req.ID,
		Name:      req.Name,
		Type:      req.Type,
		CreatorID: c.Query("userId"),
	})
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	h.hub.Rooms[room.ID] = &Room{
		ID:      room.ID,
		Name:    room.Name,
		Clients: make(map[string]*Client),
	}

	conn.WriteJSON(room)
}

func (h *WSHandler) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	roomID := c.Param("roomId")
	clientID := c.Query("userId")
	username := c.Query("username")

	// Verify room exists in database
	_, err = h.service.GetChatRoomByID(c.Request.Context(), roomID)
	if err != nil {
		conn.Close()
		return
	}

	// Store join event in database
	err = h.service.AddUserToChatRoom(c.Request.Context(), &services.AddUserToChatRoomReq{
		UserID:     clientID,
		ChatRoomID: roomID,
	})
	if err != nil {
		conn.Close()
		return
	}

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	h.hub.Register <- cl

	m := &Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	h.service.CreateMessage(c.Request.Context(), &services.CreateMessageReq{
		Content:  m.Content,
		RoomID:   m.RoomID,
		Username: m.Username,
	})

	h.hub.Broadcast <- m

	go cl.writeMessage()
	cl.readMessage(h.hub)
}

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *WSHandler) GetRoomClients(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	var clients []ClientRes
	roomId := c.Param("roomId")

	room, err := h.service.GetChatRoomByID(c.Request.Context(), roomId)
	if err != nil {
		conn.WriteJSON(gin.H{"error": "Room not found"})
		return
	}

	if _, ok := h.hub.Rooms[room.ID]; !ok {
		clients = make([]ClientRes, 0)
		conn.WriteJSON(clients)
		return
	}

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

	var req services.UpdateChatRoomReq
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
