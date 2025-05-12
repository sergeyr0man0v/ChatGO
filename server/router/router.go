package router

import (
	"chatgo/server/internal/transport"

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
	r.GET("/ws/getMessages/:roomId", wsHandler.GetMessagesByRoomID)
	r.GET("/ws/getRooms", wsHandler.GetChatRoomsByUserID)
	r.PUT("/ws/updateRoom", wsHandler.UpdateChatRoom)
	r.DELETE("/ws/deleteRoom/:roomId", wsHandler.DeleteChatRoom)

	r.POST("/ws/createRoom", wsHandler.CreateRoom)
	r.GET("/ws/joinRoom/:roomId", wsHandler.JoinRoom)
	r.GET("/ws/getAllRooms", wsHandler.GetAllRooms)
	r.GET("/ws/getRoomClients/:roomId", wsHandler.GetRoomClients)
}

func Start(addr string) error {
	return r.Run(addr)
}
