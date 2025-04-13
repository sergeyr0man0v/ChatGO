package main

import (
	"bufio"
	"bytes"
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
	ID   string `json:"id"`
	Name string `json:"name"`
}

func roomExists(serverAddr string, roomID string) bool {
	resp, err := http.Get(fmt.Sprintf("%s/ws/getRooms", serverAddr))
	if err != nil {
		log.Printf("Failed to get rooms: %v", err)
		return false
	}
	defer resp.Body.Close()

	var rooms []Room
	if err := json.NewDecoder(resp.Body).Decode(&rooms); err != nil {
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

func main() {

	serverAddr := flag.String("server", "http://localhost:8080", "Server address")
	username := flag.String("username", "", "Your username")
	password := flag.String("password", "", "Your password")
	roomID := flag.String("room", "default", "Room ID to join")
	flag.Parse()

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
		log.Printf("Signup failed: %s", string(body))
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

	if !roomExists(*serverAddr, *roomID) {
		roomData := Room{
			ID:   *roomID,
			Name: *roomID,
		}
		roomJSON, _ := json.Marshal(roomData)
		resp, err = http.Post(*serverAddr+"/ws/createRoom", "application/json", bytes.NewBuffer(roomJSON))
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

	wsURL := fmt.Sprintf("ws://localhost:8080/ws/joinRoom/%s?userId=%s&username=%s", *roomID, *username, *username)
	log.Printf("Connecting to WebSocket server at: %s", wsURL)

	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("Failed to connect to WebSocket server:", err)
	}
	defer c.Close()

	log.Printf("Successfully connected to room: %s as user: %s", *roomID, *username)

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}
			fmt.Printf("\n%s\n", string(message))
		}
	}()

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

		log.Printf("Sending message: %s", text)
		err = c.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			log.Println("Error sending message:", err)
			break
		}
	}
}
