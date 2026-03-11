package user

import (
	"github.com/google/uuid"
)

// question: should Delete() just use the userID or should it require the entire user
type UserRepository interface {
	Create(usr *User) error
	Save(usr *User) error
	Delete(usr *User) error

	// question: [userID string] or [id string]
	Get(id uuid.UUID) (*User, error)
	GetAll() ([]*User, error)
}
