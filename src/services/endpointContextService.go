package services

import (
	"orkidslearning/src/database"
)

// ContextService is a service that provides a context
type ContextService struct {
	db         *database.Database
	jwtService *JWTService
	postgres   *database.PostgresDatabase
}

// NewContextService creates a new ContextService
func NewContextService(db *database.Database, jwtService *JWTService, postgres *database.PostgresDatabase) *ContextService {
	return &ContextService{db: db, jwtService: jwtService, postgres: postgres}
}

// GetDB returns the database
func (s *ContextService) GetDB() *database.Database {
	return s.db
}

// GetJWTService returns the JWT service
func (s *ContextService) GetJWTService() *JWTService {
	return s.jwtService
}

// GetSQLDB returns the SQL database
func (s *ContextService) GetPostgres() *database.PostgresDatabase {
	return s.postgres
}
