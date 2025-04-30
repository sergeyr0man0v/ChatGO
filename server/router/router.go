package router

import (
	"server/internal/transport"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(
	userHandler *transport.UserHandler,
	wsHandler *transport.WSHandler,
) {
	r = gin.Default()

	/*r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))*/
	// User routes
	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)
	r.GET("/logout", userHandler.Logout)
	r.GET("/users", userHandler.GetAllUsers)

	// WebSocket routes
	r.GET("/rooms/:roomId/messages", wsHandler.GetMessagesByRoomID)
	r.GET("/rooms/:userId/rooms", wsHandler.GetChatRoomsByUserID)
	r.PUT("/rooms/:roomId", wsHandler.UpdateChatRoom)
	r.DELETE("/rooms/:roomId", wsHandler.DeleteChatRoom)

	r.POST("/ws/createRoom", wsHandler.CreateRoom)
	r.GET("/ws/joinRoom/:roomId", wsHandler.JoinRoom)
	r.GET("/ws/getRooms", wsHandler.GetRooms)
	r.GET("/ws/getRoomClients/:roomId", wsHandler.GetRoomClients)
}

func Start(addr string) error {
	return r.Run(addr)
}
