package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"orkidslearning/src/config"
	"orkidslearning/src/database"
	routerFunctions "orkidslearning/src/router"
	"time"

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

	// Define routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Gin server with MongoDB!"})
	})

	router.GET("/courses", routerFunctions.GetAllCourses)
	router.GET("/courses/:id", routerFunctions.GetCourseById)

	// Start server
	fmt.Printf("Server running at http://localhost:%s\n", env.Port)
	err = router.Run(":" + env.Port)
	if err != nil {
		return
	}
}
