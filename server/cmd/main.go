package main

import (
	"log"
	"server/internal/db"
	"server/internal/services"
	"server/internal/transport"
	"server/router"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Could not initialize database connection: %s", err)
	}

	userRep := db.NewRepository(dbConn.GetDB())
	userSvc := services.NewService(userRep)
	userHandler := services.NewHandler(userSvc)

	hub := transport.NewHub()
	wsHandler := transport.NewHandler(hub)
	go hub.Run()

	router.InitRouter(userHandler, wsHandler)
	router.Start("0.0.0.0:8080")
}
