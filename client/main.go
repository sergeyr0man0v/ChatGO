package main

import (
	"bufio"
	"bytes"
	"chatgo/client/color"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Room struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	CreatorID string `json:"creator_id"`
}

type LoginResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type ClientRes struct {
	ID string `json:"id"`
}

type Message struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	RoomID    string `json:"roomId"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdAt"`
}

func roomExists(serverAddr string, roomID string) bool {
	wsScheme := "ws"
	wsHost := strings.Replace(strings.Replace(serverAddr, "http://", "", 1), "https://", "", 1)
	wsURL := fmt.Sprintf("%s://%s/ws/getAllRooms", wsScheme, wsHost)
	log.Printf("Connecting to WebSocket server at: %s", wsURL)

	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Printf("Failed to connect to WebSocket server: %v", err)
		return false
	}
	defer c.Close()

	var rooms []Room
	err = c.ReadJSON(&rooms)
	if err != nil {
		log.Printf("Failed to decode rooms: %v", err)
		return false
	}

	for _, room := range rooms {
		if room.ID == roomID {
			return true
		}
	}
	return false
}

func viewAllRooms(serverAddr string) {
	wsScheme := "ws"
	wsHost := strings.Replace(strings.Replace(serverAddr, "http://", "", 1), "https://", "", 1)
	wsURL := fmt.Sprintf("%s://%s/ws/getAllRooms", wsScheme, wsHost)
	log.Printf("Connecting to WebSocket server at: %s", wsURL)

	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer c.Close()

	var rooms []Room
	err = c.ReadJSON(&rooms)
	if err != nil {
		log.Fatalf("Failed to decode rooms: %v", err)
	}

	fmt.Println("Rooms:")
	for _, room := range rooms {
		fmt.Printf("ID: %s, Name: %s, Type: %s, CreatorID: %s\n",
			room.ID, room.Name, room.Type, room.CreatorID)
	}
}

func createNewRoom(serverAddr string, roomID, roomName, roomType string, creatorID string) {
	roomData := Room{
		ID:        roomID,
		Name:      roomName,
		Type:      roomType,
		CreatorID: creatorID,
	}
	roomJSON, _ := json.Marshal(roomData)
	resp, err := http.Post(fmt.Sprintf("%s/ws/createRoom?userId=%s", serverAddr, creatorID), "application/json", bytes.NewBuffer(roomJSON))
	if err != nil {
		log.Fatalf("Failed to create room: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Room creation failed: %s", string(body))
	}

	fmt.Println("Room created successfully.")
}

func displayChatHistory(serverAddr string, roomID string, limit int, userID string) {
	wsScheme := "ws"
	wsHost := strings.Replace(strings.Replace(serverAddr, "http://", "", 1), "https://", "", 1)
	wsURL := fmt.Sprintf("%s://%s/ws/getMessages/%s/%d?userId=%s", wsScheme, wsHost, roomID, limit, userID)

	header := http.Header{}
	header.Add("Origin", serverAddr)
	header.Add("User-Agent", "ChatGO-Client")

	historyConn, _, err := websocket.DefaultDialer.Dial(wsURL, header)
	if err != nil {
		log.Printf("Failed to fetch chat history: %v", err)
		return
	}
	defer historyConn.Close()

	// Read the response
	var response interface{}
	err = historyConn.ReadJSON(&response)
	if err != nil {
		log.Printf("Failed to read response: %v", err)
		return
	}

	// Check if it's an error response
	if errorMap, ok := response.(map[string]interface{}); ok {
		if errorMsg, ok := errorMap["error"].(string); ok {
			log.Printf("Error from server: %s", errorMsg)
			return
		}
	}

	// If not an error, it should be an array of messages
	messages, ok := response.([]interface{})
	if !ok {
		log.Printf("Unexpected response format from server")
		return
	}

	fmt.Printf("\nLast %d messages:\n", limit)
	fmt.Println("----------------------------------------")
	for _, msg := range messages {
		if msgMap, ok := msg.(map[string]interface{}); ok {
			username, _ := msgMap["username"].(string)
			content, _ := msgMap["content"].(string)
			fmt.Println(color.ColorizeMessage(username, content))
		}
	}
	fmt.Println("----------------------------------------")
}

func handleMessages(c *websocket.Conn, username string, roomID string) {
	for {
		var message Message
		err := c.ReadJSON(&message)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}
		fmt.Println(color.ColorizeMessage(message.Username, message.Content))
	}
}

func viewChatHistory(serverAddr string, roomID string, limit int) {
	wsScheme := "ws"
	wsHost := strings.Replace(strings.Replace(serverAddr, "http://", "", 1), "https://", "", 1)
	wsURL := fmt.Sprintf("%s://%s/ws/getMessages/%s/%d", wsScheme, wsHost, roomID, limit)
	log.Printf("Connecting to WebSocket server at: %s", wsURL)

	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer c.Close()

	var messages []Message
	err = c.ReadJSON(&messages)
	if err != nil {
		log.Fatalf("Failed to decode messages: %v", err)
	}

	fmt.Printf("\nChat History for Room %s:\n", roomID)
	fmt.Println("----------------------------------------")
	for _, msg := range messages {
		fmt.Println(color.ColorizeMessage(msg.Username, msg.Content))
	}
	fmt.Println("----------------------------------------")
}

func main() {
	serverAddr := flag.String("server", "http://localhost:8080", "Server address")
	username := flag.String("username", "", "Username")
	password := flag.String("password", "", "Password")
	roomID := flag.String("room", "default", "Room ID to join (default: 'default')")
	createRoom := flag.Bool("createRoom", false, "Create a new room")
	roomName := flag.String("roomName", "", "Name for the new room")
	roomType := flag.String("roomType", "group", "Type of the new room")
	viewRooms := flag.Bool("viewRooms", false, "View all available rooms")
	viewHistory := flag.Bool("history", false, "View chat history")
	historyLimit := flag.Int("limit", 50, "Number of messages to retrieve for history")
	flag.Parse()

	if *username == "" || *password == "" {
		log.Fatal("Username and password are required")
	}

	// Login
	loginData := User{
		Username: *username,
		Password: *password,
	}
	loginJSON, _ := json.Marshal(loginData)
	resp, err := http.Post(fmt.Sprintf("%s/login", *serverAddr), "application/json", bytes.NewBuffer(loginJSON))
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Login failed: %s", string(body))
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		log.Fatalf("Failed to decode login response: %v", err)
	}

	if *viewRooms {
		viewAllRooms(*serverAddr)
		return
	}

	if *viewHistory {
		if *roomID == "" {
			*roomID = "default"
		}
		viewChatHistory(*serverAddr, *roomID, *historyLimit)
		return
	}

	if *createRoom {
		if *roomName == "" {
			log.Fatal("Room name is required to create a room")
		}
		createNewRoom(*serverAddr, "", *roomName, *roomType, loginResp.ID)
		return
	}

	// Check if user is already a member of the room
	wsScheme := "ws"
	wsHost := strings.Replace(strings.Replace(*serverAddr, "http://", "", 1), "https://", "", 1)
	wsURL := fmt.Sprintf("%s://%s/ws/getRoomClients/%s", wsScheme, wsHost, *roomID)
	log.Printf("Checking room members at: %s", wsURL)

	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("Failed to check room members:", err)
	}
	defer c.Close()

	var clients []ClientRes
	if err := c.ReadJSON(&clients); err != nil {
		log.Fatal("Failed to decode room members:", err)
	}

	isMember := false
	for _, client := range clients {
		if client.ID == loginResp.ID {
			isMember = true
			break
		}
	}

	if !isMember {
		log.Printf("User is not a member of room %s, joining...", *roomID)
	}

	wsScheme = "ws"
	wsHost = strings.Replace(strings.Replace(*serverAddr, "http://", "", 1), "https://", "", 1)
	wsURL = fmt.Sprintf("%s://%s/ws/joinRoom/%s?userId=%s&username=%s", wsScheme, wsHost, *roomID, loginResp.ID, *username)
	log.Printf("Connecting to WebSocket server at: %s", wsURL)

	header := http.Header{}
	header.Add("Origin", *serverAddr)
	header.Add("User-Agent", "ChatGO-Client")

	log.Printf("Attempting to connect to WebSocket server at: %s", wsURL)
	c, wsResp, err := websocket.DefaultDialer.Dial(wsURL, header)
	if err != nil {
		if wsResp != nil {
			body, _ := io.ReadAll(wsResp.Body)
			log.Fatalf("Failed to connect to WebSocket server: %v\nResponse status: %d\nResponse body: %s", err, wsResp.StatusCode, string(body))
		} else {
			log.Fatalf("Failed to connect to WebSocket server: %v", err)
		}
	}
	defer c.Close()

	log.Printf("Successfully connected to room: %s as user: %s", *roomID, *username)

	go handleMessages(c, *username, *roomID)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Connected to chat room. Type your messages (or 'exit' to quit):")
	fmt.Println("Commands:")
	fmt.Println("  /history [number] - Show last N messages (default: 10)")
	fmt.Println("  exit - Leave the chat room")

	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading input:", err)
			continue
		}

		text = strings.TrimSpace(text)
		if text == "exit" {
			break
		}

		// Handle /history command
		if strings.HasPrefix(text, "/history") {
			parts := strings.Fields(text)
			limit := 10 // default limit
			if len(parts) > 1 {
				if n, err := strconv.Atoi(parts[1]); err == nil && n > 0 {
					limit = n
				}
			}
			displayChatHistory(*serverAddr, *roomID, limit, loginResp.ID)
			continue
		}

		message := Message{
			Content:  text,
			RoomID:   *roomID,
			Username: *username,
		}

		log.Printf("Sending message to room %s: %s", *roomID, text)
		err = c.WriteJSON(message)
		if err != nil {
			log.Printf("Error sending message: %v", err)
			break
		}
		log.Printf("Message sent successfully")
	}
}
