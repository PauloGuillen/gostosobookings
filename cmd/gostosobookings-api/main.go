package main

import (
	"log"
	"net/http"

	"github.com/PauloGuillen/gostosobookings/internal/config"
	"github.com/PauloGuillen/gostosobookings/pkg/router"
)

func main() {
	// Load configuration
	config.Load()

	// Initialize DB connection
	config.InitDB()
	defer config.CloseDB()

	// Start API
	log.Println("Starting GostosoBookings API...")

	// Setup router
	r := router.SetupRouter()

	// Get server port
	port := config.Config.ServerPort
	log.Printf("Starting GostosoBookings API on port %s...", port)

	// Start the server
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
