package user

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

type SaveUserParams struct {
	name   string
	authId string
	status Status
}

type UpdateUserProfile struct {
	name      string
	avatarUrl string
	bio       string
	websites  string
	social    SocialLinks
}

func (r *Repository) CreateUser(params *SaveUserParams) (string, error) {
	user := &User{
		ID:     uuid.New(),
		Name:   strings.ToLower(params.name),
		AuthId: params.authId,
		Status: params.status,
	}

	err := r.db.Create(&user).Error

	if err != nil {
		return "", err
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

func (r *Repository) GetUserByAuthId(authId string) (*User, error) {
	var user User

	err := r.db.Where("auth_id = ?", authId).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetProfileByAuthId(authId string) (*Profile, error) {
	var u User

	err := r.db.
		Preload("Roles").
		Preload("UserProfile").
		Preload("ProducerProfile").
		Where("auth_id = ?", authId).
		First(&u).Error

	if err != nil {
		return nil, err
	}

	return &Profile{
		User:            &u,
		UserProfile:     u.UserProfile,
		ProducerProfile: u.ProducerProfile,
	}, nil
}

func (r *Repository) UpdateProfileByAuthId(params *UpdateUserProfile, authId string) (string, error) {
	var user User

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Preload("UserProfile").Preload("ProducerProfile").
			Where("auth_id = ?", authId).First(&user).Error; err != nil {
			return err
		}

		user.Name = strings.ToLower(params.name)

		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		if user.UserProfile != nil {
			user.UserProfile.AvatarUrl = params.avatarUrl
			if err := tx.Save(user.UserProfile).Error; err != nil {
				return err
			}
		}

		if user.ProducerProfile != nil {
			user.ProducerProfile.Bio = params.bio
			user.ProducerProfile.Websites = params.websites
			user.ProducerProfile.Social = params.social
			if err := tx.Save(user.ProducerProfile).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return user.ID.String(), nil
}

func (r *Repository) AddRole(userID uuid.UUID, roleID RoleType) error {
	user := &User{
		ID: userID,
	}

	err := r.db.Find(&user)

	if err != nil {
		return err.Error
	}

	return r.db.Model(user).Association("Roles").Append(&Role{RoleId: roleID})
}
