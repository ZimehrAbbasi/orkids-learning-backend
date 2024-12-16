package api

import (
	"log"
	"net/http"
	"strings"

	"orkidslearning/src/services"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware checks the validity of the token
func JWTAuthMiddleware(jwtService *services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Println("Unauthorized 1")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwtService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			log.Println("Unauthorized 2")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
func InjectContextService(contextService *services.ContextService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("contextService", contextService)
		c.Next()
	}
}
