package response

import (
	models "orkidslearning/src/models/database"
)

type AuthResponse struct {
	Message string      `json:"message"`
	User    models.User `json:"user"`
	Token   string      `json:"token"`
	Error   string      `json:"error"`
}
