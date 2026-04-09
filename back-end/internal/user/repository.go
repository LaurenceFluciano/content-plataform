package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

/* --- CREATE USER --- */

type SaveUserParams struct {
	name   string
	authId string
	status Status
}

func (r *Repository) CreateUser(params *SaveUserParams) (string, error) {
	user := &User{
		ID:     uuid.New(),
		Name:   params.name,
		AuthId: params.authId,
		Status: params.status,
	}

	result := r.db.Create(&user)

	if result.Error != nil {
		return "", result.Error
	}

	return user.ID.String(), nil
}

func (r *Repository) GetUserById(id uuid.UUID) (error, *User) {
	user := &User{ID: id}

	result := r.db.First(user)

	if result.Error != nil {
		return result.Error, nil
	}

	return nil, user
}
