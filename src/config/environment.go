package config

import (
	"os"

	"orkidslearning/src/utils/errors"

	"github.com/joho/godotenv"
)

// Environment holds all environment variables for the application
type Environment struct {
	MongoURI               string
	DBName                 string
	Port                   string
	JWTSecretKey           string
	JWTExpirationTime      string
	OTELResourceAttributes string
	FrontendURL            string
}

// LoadEnv loads environment variables into the Environment struct
func LoadEnv() (*Environment, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return nil, errors.ErrNoEnvFile
	}

	// Populate the Environment struct
	env := &Environment{
		MongoURI:               getEnv("MONGO_URI", ""),
		Port:                   getEnv("PORT", "8080"), // Default to "8080" if PORT is not set
		DBName:                 getEnv("DB_NAME", "orkidslearning"),
		JWTSecretKey:           getEnv("JWT_SECRET_KEY", ""),
		JWTExpirationTime:      getEnv("JWT_EXPIRATION_TIME", "1h"),
		OTELResourceAttributes: getEnv("OTEL_RESOURCE_ATTRIBUTES", "service.name=orkidslearning,service.version=0.1.0"),
		FrontendURL:            getEnv("FRONTEND_URL", "http://localhost:3001"),
	}

	// Validate critical environment variables
	if env.MongoURI == "" {
		return nil, errors.EnvVariableNotSet("MONGO_URI")
	}

	if env.DBName == "" {
		return nil, errors.EnvVariableNotSet("DB_NAME")
	}

	if env.JWTSecretKey == "" {
		return nil, errors.EnvVariableNotSet("JWT_SECRET_KEY")
	}

	if env.JWTExpirationTime == "" {
		return nil, errors.EnvVariableNotSet("JWT_EXPIRATION_TIME")
	}

	if env.OTELResourceAttributes == "" {
		return nil, errors.EnvVariableNotSet("OTEL_RESOURCE_ATTRIBUTES")
	}

	if env.FrontendURL == "" {
		return nil, errors.EnvVariableNotSet("FRONTEND_URL")
	}

	return env, nil
}

// getEnv retrieves an environment variable or a default value if not set
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
