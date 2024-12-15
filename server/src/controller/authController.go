package controller

import (
	"fmt"
	"log"
	"orkidslearning/src/database"
	"orkidslearning/src/models"

	"golang.org/x/crypto/bcrypt"
)

func Signup(user models.AddUser) (*models.User, error) {
	// Check if the username or email is already in use
	err := database.CheckIfUserExists(user.Username, user.Email)
	if err != nil {
		return nil, err // Return the error to the router for appropriate handling
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	user.Password = string(hashedPassword)

	// Add the user to the database
	addedUser, err := database.AddUser(user)
	if err != nil {
		return nil, err
	}

	return addedUser, nil
}

func Login(userCredentials models.LoginUser) (*models.User, error) {
	// Retrieve the user by email
	user, err := database.GetUserByEmail(userCredentials.Email)
	if err != nil {
		log.Println("Error getting user by email: ", err)
		return nil, fmt.Errorf("user not found")
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userCredentials.Password))
	if err != nil {
		log.Println("Invalid password")
		return nil, fmt.Errorf("invalid password")
	}

	user.Password = ""

	return user, nil
}
