package database

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/laurencefluciano/content-api/internal/app/exception"
	"github.com/laurencefluciano/content-api/internal/domain/user"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) user.Repository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) toDomain(table *UserTable) *user.User {
	if table == nil {
		return nil
	}

	u := user.RestoreUser(
		table.ID,
		table.AuthId,
		table.Name,
		table.Status,
	)

	if table.UserProfile != nil {
		u.RestoreProfile(table.UserProfile.AvatarUrl)
	}

	if table.ProducerProfile != nil {
		u.RestoreProducer(
			table.ProducerProfile.Bio,
			table.ProducerProfile.Websites,
			table.ProducerProfile.Social,
		)
	}

	return u
}

func (r *gormUserRepository) Save(entity *user.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if entity.IsNew() {
			return r.createFullUser(tx, entity)
		}

		return r.updateDirtyFields(tx, entity)
	})
}
func (r *gormUserRepository) FindById(id uuid.UUID) (*user.User, error) {
	var table UserTable

	err := r.db.Preload("Roles").
		Preload("UserProfile").
		Preload("ProducerProfile").
		First(&table, "id = ?", id).Error

	if err != nil {
		return nil, r.handleDBError(err)
	}

	return r.toDomain(&table), nil
}

func (r *gormUserRepository) FindByName(name string) (*user.User, error) {
	var table UserTable
	err := r.db.Preload("Roles").
		Preload("UserProfile").
		Preload("ProducerProfile").
		First(&table, "name = ?", name).Error

	if err != nil {
		return nil, r.handleDBError(err)
	}

	return r.toDomain(&table), nil
}

func (r *gormUserRepository) FindByAuthId(authId string) (*user.User, error) {
	var table UserTable
	err := r.db.Preload("Roles").
		Preload("UserProfile").
		Preload("ProducerProfile").
		First(&table, "auth_id = ?", authId).Error

	if err != nil {
		return nil, r.handleDBError(err)
	}

	return r.toDomain(&table), nil
}

func (r *gormUserRepository) updateDirtyFields(tx *gorm.DB, entity *user.User) error {
	userId := entity.ID()

	userChanges := make(map[string]any)
	userProfileChanges := make(map[string]any)
	producerProfileChanges := make(map[string]any)

	for field := range entity.GetDirtyFields() {
		switch field {
		case user.NameField:
			userChanges["name"] = entity.Name()
		case user.StatusField:
			userChanges["status"] = entity.Status()
		case user.AvatarUrlField:
			userProfileChanges["avatar_url"] = entity.Profile().AvatarUrl()
		case user.BioField:
			producerProfileChanges["bio"] = entity.Profile().Bio()
		case user.WebsitesField:
			producerProfileChanges["websites"] = pq.StringArray(entity.Profile().Websites())
		case user.SocialField:
			producerProfileChanges["social"] = entity.Profile().Social()
		case user.RolesField:
			if err := tx.Model(&UserTable{ID: userId}).Association("Roles").Replace(toRoleTable(entity.Roles())); err != nil {
				return r.handleDBError(err)
			}
		}
	}

	if len(userChanges) > 0 {
		if err := tx.Model(&UserTable{}).Where("id = ?", userId).Updates(userChanges).Error; err != nil {
			return r.handleDBError(err)
		}
	}

	if len(userProfileChanges) > 0 {
		if err := tx.Model(&UserProfileTable{}).Where("user_id = ?", userId).Updates(userProfileChanges).Error; err != nil {
			return r.handleDBError(err)
		}
	}

	if len(producerProfileChanges) > 0 {
		if err := tx.Model(&ProducerProfileTable{}).Where("user_id = ?", userId).Updates(producerProfileChanges).Error; err != nil {
			return r.handleDBError(err)
		}

	}

	entity.SetPersisted()

	return nil
}

func (r *gormUserRepository) createFullUser(tx *gorm.DB, entity *user.User) error {
	userTable := UserTable{
		ID:     entity.ID(),
		AuthId: entity.AuthId(),
		Name:   entity.Name(),
		Status: entity.Status(),
		UserProfile: &UserProfileTable{
			ID:        entity.ID(),
			AvatarUrl: "",
		},
	}

	if err := tx.Create(&userTable).Error; err != nil {
		return r.handleDBError(err)
	}

	entity.SetPersisted()
	return nil
}

func (r *gormUserRepository) CreateProducerProfile(userId uuid.UUID) error {
	profile := ProducerProfileTable{
		ID: userId,
	}

	err := r.db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&profile).Error

	if err != nil {
		return r.handleDBError(err)
	}

	return nil
}

func toRoleTable(roles []user.RoleType) []RoleTable {
	tableRoles := make([]RoleTable, len(roles))
	for i, r := range roles {
		tableRoles[i] = RoleTable{
			RoleId: r,
		}
	}
	return tableRoles
}

func (r *gormUserRepository) handleDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrNotFound
	}

	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return exception.ErrEntityConflict
		case "23503":
			return exception.ErrInvalidEntityID
		}
	}

	return fmt.Errorf("database error: %w", err)
}
