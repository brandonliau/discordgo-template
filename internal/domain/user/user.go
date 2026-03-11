package user

import (
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID
}

// todo: check that uuid doesn't exist before creating?
func NewUser() *User {
	return &User{
		ID: uuid.New(),
	}
}
