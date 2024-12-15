package router

import (
	"log"
	"net/http"
	"orkidslearning/src/controller"
	"orkidslearning/src/models"
	"orkidslearning/src/services"

	"github.com/gin-gonic/gin"
)

func SignupHandler(c *gin.Context, jwtService *services.JWTService) {
	// Parse input
	var user models.AddUser
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Error binding JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	addedUser, err := controller.Signup(user)
	if err != nil {
		log.Println("Error adding user: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user"})
		return
	}

	token, err := jwtService.GenerateToken(addedUser.Username)
	if err != nil {
		log.Println("Error generating token: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User added successfully", "user": addedUser, "token": token})
}

func LoginHandler(c *gin.Context, jwtService *services.JWTService) {

	// Parse input
	var user models.LoginUser
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Error binding JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	LoginUser, err := controller.Login(user)
	if err != nil {
		log.Println("Error logging in user: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login user"})
		return
	}

	token, err := jwtService.GenerateToken(LoginUser.Username)
	if err != nil {
		log.Println("Error generating token: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "user": LoginUser, "token": token})
}
