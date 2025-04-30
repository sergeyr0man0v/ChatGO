package main

import (
	"log"
	"server/internal/db"
	"server/internal/services"
	"server/internal/transport"
	"server/router"
)

func main() {
	// TODO: move to config file
	dbConfig := &db.Config{
		Host:     "localhost",
		Port:     "5433",
		User:     "postgres",
		Password: "password",
		DBName:   "chat-go",
		SSLMode:  "disable",
	}
	dbConn, err := db.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Could not initialize database connection: %s", err)
	}

	// Initialize repository
	repository := db.NewRepository(dbConn.GetDB())

	// Initialize service
	service := services.NewService(repository)

	// Initialize handlers
	userHandler := transport.NewUserHandler(service)

	// Initialize WebSocket hub and handler
	hub := transport.NewHub()
	wsHandler := transport.NewWSHandler(hub, &service)
	go hub.Run()

	// Initialize router with all handlers
	router.InitRouter(userHandler, wsHandler)
	router.Start("0.0.0.0:8080")
}
