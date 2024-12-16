package controller

import (
	"context"
	"fmt"
	"log"
	models "orkidslearning/src/models/database"
	services "orkidslearning/src/services"

	"go.opentelemetry.io/otel"
	"golang.org/x/crypto/bcrypt"
)

func Signup(ctx context.Context, contextService *services.ContextService, user models.AddUser) (*models.UserPostgres, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "Signup")
	defer span.End()

	// Check if the username or email is already in use
	exists, err := contextService.GetPostgres().DoesUserExist(ctx, user.Username, user.Email)
	if err != nil {
		log.Println("User with username or email already exists", err)
		return nil, err // Return the error to the router for appropriate handling
	}
	if exists {
		return nil, fmt.Errorf("user with username or email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to hash password: ", err)
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	user.Password = string(hashedPassword)

	// Add the user to the database
	addedUser, err := contextService.GetPostgres().AddUser(ctx, user)
	if err != nil {
		log.Println("Error adding user: ", err)
		return nil, err
	}

	return addedUser, nil
}

func Login(ctx context.Context, contextService *services.ContextService, userCredentials models.LoginUser) (*models.UserPostgres, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "Login")
	defer span.End()

	// Retrieve the user by email
	user, err := contextService.GetPostgres().GetUserByEmail(ctx, userCredentials.Email)
	if err != nil {
		log.Println("Error getting user by email: ", err)
		return nil, fmt.Errorf("user not found")
	}

	log.Println("User: ", userCredentials)
	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userCredentials.Password))
	if err != nil {
		log.Println("Invalid password", err)
		return nil, fmt.Errorf("invalid password")
	}

	user.Password = ""

	return user, nil
}
