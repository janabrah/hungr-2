package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	UUID      uuid.UUID `json:"uuid"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type UserResponse struct {
	Success bool `json:"success"`
	User    User `json:"user"`
}

type CreateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
