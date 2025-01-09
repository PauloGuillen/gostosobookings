package config

import "log"

type AppConfig struct {
	ServerPort string
}

var App *AppConfig

// LoadAppConfig initializes the application configuration.
func LoadAppConfig() {
	App = &AppConfig{
		ServerPort: GetEnv("SERVER_PORT", "8080"),
	}

	log.Printf("Application config loaded: %+v", App)
}
