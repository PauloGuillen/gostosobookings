package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

// Global variable to hold the app configuration
var Config *AppConfig

// Load initializes the application configuration from environment variables or .env file.
func Load() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, loading configuration from environment variables")
	}

	Config = &AppConfig{
		DBHost:     getEnv("DB_HOST", "postgres"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "gostoso_user"),
		DBPassword: getEnv("DB_PASSWORD", "gostoso_password"),
		DBName:     getEnv("DB_NAME", "gostosobookings_db"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
