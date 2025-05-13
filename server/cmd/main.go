package main

import (
	"log"

	"chatgo/server/internal/db"
	"chatgo/server/internal/services"
	"chatgo/server/internal/transport"
	"chatgo/server/router"
)

func main() {
	dbConfig := &db.Config{
		Host:     "localhost",
		Port:     "5444",
		User:     "root",
		Password: "password",
		DBName:   "chat-go",
		SSLMode:  "disable",
	}

	database, err := db.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Could not initialize the database: %v", err)
	}
	defer database.Close()

	// Initialize repository
	repository := db.NewRepository(database.GetDB())

	// Initialize service
	service := services.NewService(repository)

	// Initialize handlers
	userHandler := transport.NewUserHandler(service)

	// Initialize WebSocket hub and handler
	hub := transport.NewHub(service)
	wsHandler := transport.NewWSHandler(hub, service)
	go hub.Run()

	// Initialize router with all handlers
	router.InitRouter(userHandler, wsHandler)
	router.Start("0.0.0.0:8080")
}
