package repository

import (
	"github.com/laurencefluciano/content-api/internal/domain/user"
	"github.com/laurencefluciano/content-api/internal/infra/database"
)

func (r *gormUserRepository) Create(entity *user.User) error {
	userTable := database.UserTable{
		ID:     entity.ID(),
		AuthId: entity.AuthId(),
		Name:   entity.Name(),
		Status: entity.Status(),
		Roles:  toRoleTable(entity.Roles()),
		UserProfile: &database.UserProfileTable{
			ID:        entity.ID(),
			AvatarUrl: "",
		},
	}

	if err := r.db.Create(&userTable).Error; err != nil {
		return r.handleDBError(err)
	}

	entity.ClearDirty()
	return nil
}
