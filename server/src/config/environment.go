package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Environment holds all environment variables for the application
type Environment struct {
	MongoURI string
	Port     string
}

// LoadEnv loads environment variables into the Environment struct
func LoadEnv() *Environment {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// Populate the Environment struct
	env := &Environment{
		MongoURI: getEnv("MONGO_URI", ""),
		Port:     getEnv("PORT", "8080"), // Default to "8080" if PORT is not set
	}

	// Validate critical environment variables
	if env.MongoURI == "" {
		log.Fatal("Environment variable MONGO_URI is required but not set")
	}

	return env
}

// getEnv retrieves an environment variable or a default value if not set
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
