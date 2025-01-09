package config

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

// LoadDatabaseConfig initializes the database connection using environment variables.
func LoadDatabaseConfig() error {
	// Ensure the environment variables are loaded
	LoadEnv()

	// Retrieve database configuration from environment variables
	dbHost := GetEnv("DB_HOST", "localhost")
	dbPort := GetEnv("DB_PORT", "5432")
	dbUser := GetEnv("DB_USER", "user")
	dbPassword := GetEnv("DB_PASSWORD", "password")
	dbName := GetEnv("DB_NAME", "database_name")

	// Build the connection URL
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	// Connect to the database
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return err
	}

	DB = conn
	log.Println("Successfully connected to the database!")
	return nil
}

// CloseDB closes the database connection gracefully.
func CloseDB() error {
	if DB != nil {
		if err := DB.Close(context.Background()); err != nil {
			log.Printf("Error closing database connection: %v", err)
			return err
		}
		log.Println("Database connection closed.")
	}
	return nil
}
