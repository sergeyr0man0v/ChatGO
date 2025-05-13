package transport

import (
	"chatgo/server/internal/interfaces"
	"context"
	"log"
)

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
	service    interfaces.Service
}

func NewHub(service interfaces.Service) *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
		service:    service,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			log.Printf("Client %s registering for room %s", cl.ID, cl.RoomID)
			if _, ok := h.Rooms[cl.RoomID]; ok {
				r := h.Rooms[cl.RoomID]

				if _, ok := r.Clients[cl.ID]; !ok {
					r.Clients[cl.ID] = cl
					log.Printf("Client %s added to room %s", cl.ID, cl.RoomID)
				}
			} else {
				log.Printf("Room %s not found in hub", cl.RoomID)
			}
		case cl := <-h.Unregister:
			log.Printf("Client %s unregistering from room %s", cl.ID, cl.RoomID)
			if _, ok := h.Rooms[cl.RoomID]; ok {
				if _, ok := h.Rooms[cl.RoomID].Clients[cl.ID]; ok {
					if len(h.Rooms[cl.RoomID].Clients) != 0 {
						h.Broadcast <- &Message{
							Content:  "User left the chat",
							RoomID:   cl.RoomID,
							Username: cl.Username,
						}
					}

					delete(h.Rooms[cl.RoomID].Clients, cl.ID)
					close(cl.Message)
					log.Printf("Client %s removed from room %s", cl.ID, cl.RoomID)
				}
			}

		case m := <-h.Broadcast:
			log.Printf("Broadcasting message to room %s: %s", m.RoomID, m.Content)
			if _, ok := h.Rooms[m.RoomID]; ok {
				// Store message in database
				_, err := h.service.CreateMessage(context.Background(), &interfaces.CreateMessageReq{
					Content:  m.Content,
					RoomID:   m.RoomID,
					Username: m.Username,
				})
				if err != nil {
					log.Printf("Failed to store message: %v", err)
				}

				for _, cl := range h.Rooms[m.RoomID].Clients {
					log.Printf("Sending message to client %s in room %s", cl.ID, m.RoomID)
					cl.Message <- m
				}
			} else {
				log.Printf("Room %s not found for message broadcast", m.RoomID)
			}
		}
	}
}
