package database

import (
	"fmt"

	"github.com/laurencefluciano/content-api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewPostgresDB() (*gorm.DB, error) {
	dsn := config.GetEnv("DATABASE_URL")
	dbSchema := config.GetEnv("DB_SCHEMA")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: fmt.Sprintf("%s.", dbSchema),
		},
	})

	if err != nil {
		return nil, err
	}

	return db, err
}
