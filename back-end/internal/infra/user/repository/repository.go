package repository

import (
	"github.com/laurencefluciano/content-api/internal/domain/user"
	"gorm.io/gorm"
)

type gormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) user.Repository {
	return &gormUserRepository{db: db}
}
