// main.go
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

	// Ensure that the database connection will be closed properly at the end
	defer config.CloseDB()

	// Start the API server
	startServer()
}

// initializeApp loads environment variables and initializes the database connection.
func initializeApp() error {
	// Load environment variables from .env file or system environment
	config.LoadEnv()

	// Initialize the database connection
	if err := config.LoadDatabaseConfig(); err != nil {
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

	// Initialize the Auth service
	authService := service.NewAuthService(userRepository)

	// Set up the router with userService and authService
	r := router.SetupRouter(*userService, *authService)

	// Retrieve and log the server port
	port := config.GetEnv("SERVER_PORT", "8080")
	log.Printf("Starting GostosoBookings API on port %s...", port)

	// Start the server
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
