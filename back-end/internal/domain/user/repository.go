package user

import (
	"github.com/google/uuid"
)

type Repository interface {
	Create(params *User) error
	Update(params *User) error
	FindById(id uuid.UUID) (*User, error)
	FindByAuthId(authId string) (*User, error)
}
