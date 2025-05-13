package main

import (
	"log"

	"chatgo/server/internal/db"
	"chatgo/server/internal/services"
	"chatgo/server/internal/transport"
	"chatgo/server/pkg/config"
	"chatgo/server/router"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("server/pkg/config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	database, err := db.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Could not initialize the database: %v", err)
	}
	defer database.Close()

	// Initialize repository
	repository := db.NewRepository(database.GetDB())

	// Initialize service
	service := services.NewService(repository, &cfg.Service)

	// Initialize handlers
	userHandler := transport.NewUserHandler(service)

	// Initialize WebSocket hub and handler
	hub := transport.NewHub(service)
	wsHandler := transport.NewWSHandler(hub, service)
	go hub.Run()

	// Initialize router with all handlers
	router.InitRouter(userHandler, wsHandler)
	router.Start(&cfg.Server)
}
