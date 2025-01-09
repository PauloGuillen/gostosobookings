package main

import (
	"log"
	"net/http"

	"github.com/PauloGuillen/gostosobookings/config"
	"github.com/PauloGuillen/gostosobookings/internal/user/repository"
	"github.com/PauloGuillen/gostosobookings/internal/user/service"
	"github.com/PauloGuillen/gostosobookings/pkg/router"
)

func main() {
	// Load environment variables and initialize database connection
	if err := initializeApp(); err != nil {
		log.Fatalf("Error initializing application: %v", err)
	}

	// Start the API server
	startServer()

	// Ensure that the database connection will be closed properly at the end
	defer config.CloseDB()
}

// initializeApp loads environment variables and initializes the database connection.
func initializeApp() error {
	// Load environment variables from .env file or system environment
	config.LoadEnv() // Calling LoadEnv() to load .env variables

	// Initialize the database connection
	if err := config.LoadDatabaseConfig(); err != nil { // Renamed to LoadDatabaseConfig
		return err
	}

	log.Println("Application initialized successfully.")
	return nil
}

// startServer configures and starts the HTTP server.
func startServer() {
	// Initialize the User repository and service
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository)

	// Set up the router with userService
	r := router.SetupRouter(*userService)

	// Retrieve and log the server port
	port := config.GetEnv("SERVER_PORT", "8080") // Use GetEnv to fetch the port
	log.Printf("Starting GostosoBookings API on port %s...", port)

	// Start the server
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
