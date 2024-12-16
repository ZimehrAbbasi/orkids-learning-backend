package main

import (
	"context"
	"fmt"
	"log"
	api "orkidslearning/src/api"
	"orkidslearning/src/config"
	"orkidslearning/src/database"
	"orkidslearning/src/services"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// Load environment variables
	var err error
	env, err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.NewDatabase(ctx, env.MongoURI, env.DBName)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Disconnect(ctx)

	var jwtService = services.NewJWTService(env.JWTSecretKey, env.JWTExpirationTime)
	var contextService = services.NewContextService(db, jwtService)

	// Create a Gin router
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Define protected router group with JWT middleware
	api.InitializeRoutes(router, contextService)

	// Start server
	fmt.Printf("Server running at http://localhost:%s\n", env.Port)
	err = router.Run(":" + env.Port)
	if err != nil {
		log.Fatal(err)
		return
	}
}
