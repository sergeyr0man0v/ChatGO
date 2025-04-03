package main
// server
import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	rdb *redis.Client
	db  *gorm.DB
)

type Message struct {
	gorm.Model
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

func main() {
	initDB()
	initRedis()
	
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/messages", getMessages)
	
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		
		// save to DB
		message := Message{Sender: "user", Content: string(msg)}
		db.Create(&message)
		
		// publish to Redis
		rdb.Publish(r.Context(), "chat", msg)
	}
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	var messages []Message
	db.Find(&messages)
	// JSON response implementation todo
}

func initDB() {
	// PostgreSQL initialization todo
}

func initRedis() {
	// Redis client init todo
}
