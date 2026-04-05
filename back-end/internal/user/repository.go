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

func (r *Repository) CreateUser(user *User) error {

	user.ID = uuid.New()

	result := r.db.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *Repository) GetUserById(id uuid.UUID) (error, *User) {
	user := &User{ID: id}

	result := r.db.First(user)

	if result.Error != nil {
		return result.Error, nil
	}

	return nil, user
}
