package router

import (
	"context"
	"log"
	"net/http"
	"orkidslearning/src/controller"
	models "orkidslearning/src/models/database"
	"orkidslearning/src/models/response"
	"orkidslearning/src/services"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

func SignupHandler(c *gin.Context) {
	tracer := otel.Tracer("router")
	ctx, span := tracer.Start(c.Request.Context(), "SignupHandler")
	defer span.End()

	contextService, exists := c.MustGet("contextService").(*services.ContextService)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load contextService"})
		return
	}

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

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	addedUser, err := controller.Signup(ctx, contextService.GetDB(), user)
	if err != nil {
		log.Println("Error adding user: ", err)
		c.JSON(http.StatusInternalServerError, response.AuthResponse{
			Message: "Failed to add user",
			Error:   err.Error(),
		})
		return
	}

	token, err := contextService.GetJWTService().GenerateToken(addedUser.Username)
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

func LoginHandler(c *gin.Context) {
	tracer := otel.Tracer("router")
	ctx, span := tracer.Start(c.Request.Context(), "LoginHandler")
	defer span.End()

	contextService, exists := c.MustGet("contextService").(*services.ContextService)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load contextService"})
		return
	}

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

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	LoginUser, err := controller.Login(ctx, contextService.GetDB(), user)
	if err != nil {
		log.Println("Error logging in user: ", err)
		c.JSON(http.StatusInternalServerError, response.AuthResponse{
			Message: "Failed to login user",
			Error:   err.Error(),
		})
		return
	}

	token, err := contextService.GetJWTService().GenerateToken(LoginUser.Username)
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
