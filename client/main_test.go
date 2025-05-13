package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestRoomExists(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rooms := []Room{
			{ID: "1", Name: "Room 1", Type: "group", CreatorID: "user1"},
			{ID: "2", Name: "Room 2", Type: "group", CreatorID: "user2"},
		}
		json.NewEncoder(w).Encode(rooms)
	}))
	defer server.Close()

	// Test cases
	tests := []struct {
		name     string
		roomID   string
		expected bool
	}{
		{"Existing Room", "1", true},
		{"Non-existing Room", "3", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := roomExists(server.URL, tt.roomID)
			if result != tt.expected {
				t.Errorf("roomExists() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCreateNewRoom(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check if userId is present in query
		if r.URL.Query().Get("userId") == "" {
			http.Error(w, "user ID is required", http.StatusBadRequest)
			return
		}

		// Parse request body
		var room Room
		if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Return success response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(room)
	}))
	defer server.Close()

	// Test cases
	tests := []struct {
		name      string
		roomID    string
		roomName  string
		roomType  string
		creatorID string
		wantErr   bool
	}{
		{
			name:      "Valid Room Creation",
			roomID:    "",
			roomName:  "Test Room",
			roomType:  "group",
			creatorID: "user1",
			wantErr:   false,
		},
		{
			name:      "Missing Room Name",
			roomID:    "",
			roomName:  "",
			roomType:  "group",
			creatorID: "user1",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a buffer to capture output
			var buf bytes.Buffer
			log.SetOutput(&buf)

			// Call the function
			createNewRoom(server.URL, tt.roomID, tt.roomName, tt.roomType, tt.creatorID)

			// Check if error was logged
			if tt.wantErr {
				if buf.Len() == 0 {
					t.Error("Expected error to be logged, but no error was found")
				}
			} else {
				if buf.Len() > 0 {
					t.Errorf("Unexpected error: %s", buf.String())
				}
			}
		})
	}
}

func TestDisplayChatHistory(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Upgrade to WebSocket
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("Failed to upgrade connection: %v", err)
			return
		}
		defer conn.Close()

		// Send test messages
		messages := []Message{
			{
				ID:        "1",
				Content:   "Test message 1",
				RoomID:    "room1",
				Username:  "user1",
				CreatedAt: "2024-01-01T00:00:00Z",
			},
			{
				ID:        "2",
				Content:   "Test message 2",
				RoomID:    "room1",
				Username:  "user2",
				CreatedAt: "2024-01-01T00:01:00Z",
			},
		}
		conn.WriteJSON(messages)
	}))
	defer server.Close()

	// Test cases
	tests := []struct {
		name    string
		roomID  string
		limit   int
		userID  string
		wantErr bool
	}{
		{
			name:    "Valid History Request",
			roomID:  "room1",
			limit:   10,
			userID:  "user1",
			wantErr: false,
		},
		{
			name:    "Invalid Room ID",
			roomID:  "invalid",
			limit:   10,
			userID:  "user1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a buffer to capture output
			var buf bytes.Buffer
			log.SetOutput(&buf)

			// Call the function
			displayChatHistory(server.URL, tt.roomID, tt.limit, tt.userID)

			// Check if error was logged
			if tt.wantErr {
				if buf.Len() == 0 {
					t.Error("Expected error to be logged, but no error was found")
				}
			} else {
				if buf.Len() > 0 {
					t.Errorf("Unexpected error: %s", buf.String())
				}
			}
		})
	}
}

func TestHandleMessages(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Upgrade to WebSocket
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("Failed to upgrade connection: %v", err)
			return
		}
		defer conn.Close()

		// Send test messages
		messages := []Message{
			{
				ID:        "1",
				Content:   "Test message 1",
				RoomID:    "room1",
				Username:  "user1",
				CreatedAt: "2024-01-01T00:00:00Z",
			},
			{
				ID:        "2",
				Content:   "Test message 2",
				RoomID:    "room1",
				Username:  "user2",
				CreatedAt: "2024-01-01T00:01:00Z",
			},
		}

		for _, msg := range messages {
			if err := conn.WriteJSON(msg); err != nil {
				t.Errorf("Failed to send message: %v", err)
				return
			}
		}
	}))
	defer server.Close()

	// Create a WebSocket connection
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close()

	// Create a channel to receive messages
	msgChan := make(chan Message, 2)
	go func() {
		for {
			var msg Message
			if err := conn.ReadJSON(&msg); err != nil {
				return
			}
			msgChan <- msg
		}
	}()

	// Test message handling
	username := "testuser"
	roomID := "room1"
	go handleMessages(conn, username, roomID)

	// Wait for messages
	for i := 0; i < 2; i++ {
		select {
		case msg := <-msgChan:
			if msg.RoomID != roomID {
				t.Errorf("Expected room ID %s, got %s", roomID, msg.RoomID)
			}
		case <-time.After(time.Second):
			t.Error("Timeout waiting for message")
		}
	}
}
