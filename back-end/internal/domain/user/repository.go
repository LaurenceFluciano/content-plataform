package user

import (
	"github.com/google/uuid"
)

type SaveUserParams struct {
	Name   string
	AuthId string
	Status Status
}

type UpdateUserParams struct {
	Name      string
	AvatarUrl string
	Bio       string
	Websites  []string
	Social    *SocialLinks
}

type Repository interface {
	Save(params *User) error
	FindById(id uuid.UUID) (*User, error)
	FindByName(name string) (*User, error)
	FindByAuthId(authId string) (*User, error)
}
