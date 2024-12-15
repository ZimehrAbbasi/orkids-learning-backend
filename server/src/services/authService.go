package services

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey      string
	expirationTime time.Duration
}

func NewJWTService(secretKey string, expirationDurationString string) *JWTService {
	expirationDuration, err := time.ParseDuration(expirationDurationString)
	if err != nil {
		log.Fatal("Invalid JWT expiration time:", err)
	}
	return &JWTService{secretKey: secretKey, expirationTime: expirationDuration}
}

// GenerateToken creates a new JWT token
func (s *JWTService) GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(s.expirationTime).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

// ValidateToken validates a JWT token
func (s *JWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return []byte(s.secretKey), nil
	})
}

func (s *JWTService) GetUsernameFromToken(tokenString string) (string, error) {
	token, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return token.Claims.(jwt.MapClaims)["username"].(string), nil
}
