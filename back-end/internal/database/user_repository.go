package database

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/laurencefluciano/content-api/internal/exceptions"
	"github.com/laurencefluciano/content-api/internal/user"
	"gorm.io/gorm"
)

func mapToDomain(u *UserTable) *user.User {
	var roles []user.RoleType

	for _, role := range u.Roles {
		roles = append(roles, role.RoleId)
	}

	return &user.User{
		ID:     u.ID,
		AuthId: u.AuthId,
		Name:   u.Name,
		Roles:  roles,
		Profile: &user.Profile{
			Bio:       u.ProducerProfile.Bio,
			AvatarUrl: u.UserProfile.AvatarUrl,
			Social:    u.ProducerProfile.Social,
			Websites:  u.ProducerProfile.Websites,
		},
	}
}

type gormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) user.Repository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) Register(params *user.SaveUserParams) error {
	var userTable UserTable

	err := r.db.Transaction(func(tx *gorm.DB) error {
		userTable = UserTable{
			ID:     uuid.New(),
			Name:   params.Name,
			AuthId: params.AuthId,
			Status: params.Status,
		}
		if err := tx.Create(&userTable).Error; err != nil {
			return err
		}

		profile := UserProfileTable{
			ID:        userTable.ID,
			AvatarUrl: "",
		}
		if err := tx.Create(&profile).Error; err != nil {
			return err
		}

		if err := tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", userTable.ID, user.RoleReader).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || errors.Is(err, gorm.ErrCheckConstraintViolated) {
			return exceptions.ErrEntityConflict
		}
		return err
	}

	return nil
}

func (r *gormUserRepository) GetById(id uuid.UUID) (*user.User, error) {
	var u UserTable

	err := r.db.
		Preload("RolesTable").
		Preload("UserProfileTable").
		Preload("ProducerProfileTable").
		Where("id = ?", id).
		First(&u).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrNotFound
		}

		return nil, err
	}

	return mapToDomain(&u), nil
}

func (r *gormUserRepository) GetByName(name string) (*user.User, error) {
	var u UserTable

	err := r.db.
		Preload("RolesTable").
		Preload("UserProfileTable").
		Preload("ProducerProfileTable").
		Where("name = ?", name).
		First(&u).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrNotFound
		}

		return nil, err
	}

	return mapToDomain(&u), nil
}

func (r *gormUserRepository) GetByAuthId(authId string) (*user.User, error) {
	var u UserTable

	err := r.db.
		Preload("RolesTable").
		Preload("UserProfileTable").
		Preload("ProducerProfileTable").
		Where("auth_id = ?", authId).
		First(&u).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrNotFound
		}

		return nil, err
	}

	return mapToDomain(&u), nil
}

func (r *gormUserRepository) SetStatusById(userId string) error {
	var userTable UserTable

	err := r.db.Where("id = ?", userId).First(&userTable).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exceptions.ErrNotFound
		}

		return err
	}

	return nil
}

func (r *gormUserRepository) UpdateByAuthId(params *user.UpdateUserParams, authId string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var userTable UserTable

		if err := tx.Where("auth_id = ?", authId).First(&userTable).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exceptions.ErrNotFound
			}
			return err
		}

		userUpdates := make(map[string]interface{})
		if params.Name != "" {
			userUpdates["name"] = strings.ToLower(params.Name)
		}
		if len(userUpdates) > 0 {
			tx.Model(&userTable).Updates(userUpdates)
		}

		if params.AvatarUrl != "" {
			tx.Model(&UserProfileTable{}).Where("user_id = ?", userTable.ID).Update("avatar_url", params.AvatarUrl)
		}

		producerUpdates := make(map[string]interface{})
		if params.Bio != "" {
			producerUpdates["bio"] = params.Bio
		}
		if len(params.Websites) > 0 {
			producerUpdates["websites"] = params.Websites
		}

		if params.Social != nil {
			producerUpdates["social"] = params.Social
		}

		if len(producerUpdates) > 0 {
			err := tx.Model(&userTable).Association("Roles").Append(&RoleTable{RoleId: user.RoleProducer})
			if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
				return exceptions.ErrEntityConflict
			}

			var pProfile ProducerProfileTable
			err = tx.Where(ProducerProfileTable{ID: userTable.ID}).
				Assign(producerUpdates).
				FirstOrCreate(&pProfile).Error

			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *gormUserRepository) SetRole(userID uuid.UUID, roleID user.RoleType) error {
	user := &UserTable{
		ID: userID,
	}

	err := r.db.Find(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exceptions.ErrNotFound
		}

		return err
	}

	err = r.db.Model(user).Association("RolesTable").Append(&RoleTable{RoleId: roleID})

	if err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return exceptions.ErrEntityConflict
		}
		return err
	}

	return nil
}
