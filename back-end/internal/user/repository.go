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
	Register(params *SaveUserParams) error
	GetById(id uuid.UUID) (*User, error)
	GetByName(name string) (*User, error)
	GetByAuthId(authId string) (*User, error)
	SetStatusById(userId string) error
	UpdateByAuthId(params *UpdateUserParams, authId string) error
	SetRole(userID uuid.UUID, roleID RoleType) error
}
