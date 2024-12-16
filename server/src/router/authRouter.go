package router

import (
	"log"
	"net/http"
	"orkidslearning/src/controller"
	models "orkidslearning/src/models/database"
	"orkidslearning/src/models/response"
	"orkidslearning/src/services"

	"github.com/gin-gonic/gin"
)

func SignupHandler(c *gin.Context, jwtService *services.JWTService) {
	// Parse input
	var user models.AddUser
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Error binding JSON: ", err)
		c.JSON(http.StatusBadRequest, response.AuthResponse{
			Message: "Error binding JSON",
			Error:   err.Error(),
		})
		return
	}

	addedUser, err := controller.Signup(user)
	if err != nil {
		log.Println("Error adding user: ", err)
		c.JSON(http.StatusInternalServerError, response.AuthResponse{
			Message: "Failed to add user",
			Error:   err.Error(),
		})
		return
	}

	token, err := jwtService.GenerateToken(addedUser.Username)
	if err != nil {
		log.Println("Error generating token: ", err)
		c.JSON(http.StatusInternalServerError, response.AuthResponse{
			Message: "Failed to generate token",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.AuthResponse{
		Message: "User added successfully",
		User:    *addedUser,
		Token:   token,
	})
}

func LoginHandler(c *gin.Context, jwtService *services.JWTService) {

	// Parse input
	var user models.LoginUser
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Error binding JSON: ", err)
		c.JSON(http.StatusBadRequest, response.AuthResponse{
			Message: "Error binding JSON",
			Error:   err.Error(),
		})
		return
	}

	LoginUser, err := controller.Login(user)
	if err != nil {
		log.Println("Error logging in user: ", err)
		c.JSON(http.StatusInternalServerError, response.AuthResponse{
			Message: "Failed to login user",
			Error:   err.Error(),
		})
		return
	}

	token, err := jwtService.GenerateToken(LoginUser.Username)
	if err != nil {
		log.Println("Error generating token: ", err)
		c.JSON(http.StatusInternalServerError, response.AuthResponse{
			Message: "Failed to generate token",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.AuthResponse{
		Message: "User logged in successfully",
		User:    *LoginUser,
		Token:   token,
	})
}
