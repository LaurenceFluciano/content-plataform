package database

import (
	"github.com/laurencefluciano/content-api/internal/config"
	"github.com/laurencefluciano/content-api/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB() (*gorm.DB, error) {
	dsn := config.GetEnv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&user.User{},
		&user.Profile{},
		&user.Producer{},
		&user.Role{},
	)

	return db, err
}
