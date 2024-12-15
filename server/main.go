package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"orkidslearning/src/config"
	"orkidslearning/src/database"
	"orkidslearning/src/middlewares"
	routerFunctions "orkidslearning/src/router"
	"orkidslearning/src/services"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	var err error
	env := config.LoadEnv()

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = database.Connect(ctx, env.MongoURI)
	if err != nil {
		log.Fatal(err)
	}

	// Create a Gin router
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	expirationDuration, err := time.ParseDuration(env.JWTExpirationTime)
	if err != nil {
		log.Fatal("Invalid JWT expiration time:", err)
	}
	var jwtService = services.NewJWTService(env.JWTSecretKey, expirationDuration)

	// Define protected router group with JWT middleware
	protected := router.Group("/protected")
	protected.Use(middlewares.JWTAuthMiddleware(jwtService))

	// Define routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Gin server with MongoDB!"})
	})

	// Course routes
	router.GET("/api/courses", routerFunctions.GetAllCourses)
	router.GET("/api/courses/:id", routerFunctions.GetCourseById)

	// Auth routes
	router.POST("/api/auth/login", func(ctx *gin.Context) {
		routerFunctions.LoginHandler(ctx, jwtService)
	})
	router.POST("/api/auth/signup", func(ctx *gin.Context) {
		routerFunctions.SignupHandler(ctx, jwtService)
	})

	// Protected routes
	protected.POST("/api/courses", routerFunctions.AddCourse)

	// Start server
	fmt.Printf("Server running at http://localhost:%s\n", env.Port)
	err = router.Run(":" + env.Port)
	if err != nil {
		return
	}
}
