package repository

import (
	"log"
	"time"

	"github.com/laurencefluciano/content-api/internal/domain/user"
	"github.com/laurencefluciano/content-api/internal/infra/database"
)

func (r *gormUserRepository) FindByAuthId(authId string) (*user.User, error) {
	start := time.Now()
	var table database.UserTable
	err := r.db.
		Joins("UserProfile").
		Joins("ProducerProfile").
		Preload("Roles").
		First(&table, "users.auth_id = ?", authId).Error

	if err != nil {
		return nil, r.handleDBError(err)
	}

	duration := time.Since(start)
	log.Printf("🕒 [REPO SQL] FindByAuthID levou: %v", duration)
	return r.toDomain(&table), nil
}
