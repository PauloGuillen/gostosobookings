package config

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

// InitDB initializes the database connection using the loaded configuration.
func InitDB() {
	if Config == nil {
		log.Fatal("AppConfig is not initialized. Call config.Load() before InitDB().")
	}

	// Build the connection URL
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		Config.DBUser,
		Config.DBPassword,
		Config.DBHost,
		Config.DBPort,
		Config.DBName,
	)

	// Connect to the database
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	DB = conn
	log.Println("Successfully connected to the database!")
}

// CloseDB closes the database connection gracefully.
func CloseDB() {
	if DB != nil {
		if err := DB.Close(context.Background()); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed.")
		}
	}
}
