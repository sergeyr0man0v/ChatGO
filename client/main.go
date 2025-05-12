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
	Content  string `json:"content"`
	RoomID   string `json:"room_id"`
	Username string `json:"username"`
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

func createRoom(serverAddr string, roomID, roomName, roomType string, creatorID string) {
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

func handleMessages(c *websocket.Conn, username string) {
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

func main() {
	serverAddr := flag.String("server", "http://localhost:8080", "Server address")
	username := flag.String("username", "", "Your username")
	password := flag.String("password", "", "Your password")
	roomID := flag.String("room", "default", "Room ID to join")
	viewRooms := flag.Bool("viewRooms", false, "View all rooms")
	createRoomFlag := flag.Bool("createRoom", false, "Create a new room")
	roomName := flag.String("roomName", "", "Name of the room to create")
	roomType := flag.String("roomType", "group", "Type of the room to create")
	flag.Parse()

	if *viewRooms {
		viewAllRooms(*serverAddr)
		return
	}

	if *createRoomFlag {
		if *roomName == "" {
			log.Fatal("Room name is required")
		}
		loginData := User{
			Username: *username,
			Password: *password,
		}
		loginJSON, _ := json.Marshal(loginData)
		resp, err := http.Post(*serverAddr+"/login", "application/json", bytes.NewBuffer(loginJSON))
		if err != nil {
			log.Fatal("Failed to login:", err)
		}
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Fatal("Login failed:", string(body))
		}

		var loginResp LoginResponse
		if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
			log.Fatal("Failed to decode login response:", err)
		}

		createRoom(*serverAddr, *roomID, *roomName, *roomType, loginResp.ID)
		return
	}

	if *username == "" || *password == "" {
		log.Fatal("Username and password are required")
	}

	signupData := User{
		Username: *username,
		Password: *password,
	}
	signupJSON, _ := json.Marshal(signupData)
	resp, err := http.Post(*serverAddr+"/signup", "application/json", bytes.NewBuffer(signupJSON))
	if err != nil {
		log.Fatal("Failed to sign up:", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if strings.Contains(string(body), "duplicate key value") {
			log.Printf("User already exists, proceeding to login...")
		} else {
			log.Fatalf("Signup failed: %s", string(body))
		}
	}

	loginData := User{
		Username: *username,
		Password: *password,
	}
	loginJSON, _ := json.Marshal(loginData)
	resp, err = http.Post(*serverAddr+"/login", "application/json", bytes.NewBuffer(loginJSON))
	if err != nil {
		log.Fatal("Failed to login:", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatal("Login failed:", string(body))
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		log.Fatal("Failed to decode login response:", err)
	}
	creatorID := loginResp.ID

	if !roomExists(*serverAddr, *roomID) {
		roomData := Room{
			ID:        *roomID,
			Name:      *roomID,
			Type:      "group",
			CreatorID: creatorID,
		}
		roomJSON, _ := json.Marshal(roomData)
		resp, err = http.Post(fmt.Sprintf("%s/ws/createRoom?userId=%s", *serverAddr, creatorID), "application/json", bytes.NewBuffer(roomJSON))
		if err != nil {
			log.Fatal("Failed to create room:", err)
		}
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Fatal("Room creation failed:", string(body))
		}
		log.Printf("Created new room: %s", *roomID)
	} else {
		log.Printf("Room %s already exists, joining...", *roomID)
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
		if client.ID == creatorID {
			isMember = true
			break
		}
	}

	if !isMember {
		log.Printf("User is not a member of room %s, joining...", *roomID)
	}

	wsScheme = "ws"
	wsHost = strings.Replace(strings.Replace(*serverAddr, "http://", "", 1), "https://", "", 1)
	wsURL = fmt.Sprintf("%s://%s/ws/joinRoom/%s?userId=%s&username=%s", wsScheme, wsHost, *roomID, creatorID, *username)
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

	go handleMessages(c, *username)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Connected to chat room. Type your messages (or 'exit' to quit):")

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
