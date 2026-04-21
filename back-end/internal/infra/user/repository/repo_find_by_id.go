package repository

import (
	"github.com/google/uuid"
	"github.com/laurencefluciano/content-api/internal/domain/user"
	"github.com/laurencefluciano/content-api/internal/infra/database"
)

func (r *gormUserRepository) FindById(id uuid.UUID) (*user.User, error) {
	var table database.UserTable

	err := r.db.Preload("Roles").
		Preload("UserProfile").
		Preload("ProducerProfile").
		First(&table, "id = ?", id).Error

	if err != nil {
		return nil, r.handleDBError(err)
	}

	return r.toDomain(&table), nil
}
