package middlewares

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
		log.Println("Auth header: ", authHeader)
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Println("Unauthorized")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwtService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			log.Println("Unauthorized")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
