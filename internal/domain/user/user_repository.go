package user

import (
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(usr *User) error
	Save(usr *User) error
	Delete(usr *User) error

	Get(id uuid.UUID) (*User, error)
	GetAll() ([]*User, error)
}
