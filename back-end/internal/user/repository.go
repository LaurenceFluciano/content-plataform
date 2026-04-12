package user

import (
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
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
	websites  []string
	social    SocialLinks
}

func (r *Repository) RegisterNewUser(params *SaveUserParams) (*User, error) {
	var user User

	err := r.db.Transaction(func(tx *gorm.DB) error {
		user = User{
			ID:     uuid.New(),
			Name:   strings.ToLower(params.name),
			AuthId: params.authId,
			Status: params.status,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		profile := userProfile{
			UserId:    user.ID,
			AvatarUrl: "",
		}
		if err := tx.Create(&profile).Error; err != nil {
			return err
		}

		if err := tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", user.ID, RoleReader).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserById(id uuid.UUID) (error, *User) {
	user := &User{ID: id}

	result := r.db.First(user)

	if result.Error != nil {
		return result.Error, nil
	}

	return nil, user
}

func (r *Repository) GetUserByName(name string) (error, *User) {
	user := &User{Name: name}

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

func (r *Repository) UpdateUserStatusById(userId uuid.UUID, status Status) error {
	var user User

	err := r.db.Where("id = ?", userId).First(&user).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateProfileByAuthId(params *UpdateUserProfile, authId string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var user User
		if err := tx.Where("auth_id = ?", authId).First(&user).Error; err != nil {
			return err
		}

		userUpdates := make(map[string]interface{})
		if params.name != "" {
			userUpdates["name"] = strings.ToLower(params.name)
		}
		if len(userUpdates) > 0 {
			tx.Model(&user).Updates(userUpdates)
		}

		if params.avatarUrl != "" {
			tx.Model(&userProfile{}).Where("user_id = ?", user.ID).Update("avatar_url", params.avatarUrl)
		}

		producerUpdates := make(map[string]interface{})
		if params.bio != "" {
			producerUpdates["bio"] = params.bio
		}
		if len(params.websites) > 0 {
			producerUpdates["websites"] = params.websites
		}

		if len(producerUpdates) > 0 {
			tx.Model(&producerProfile{}).Where("user_id = ?", user.ID).Updates(producerUpdates)
		}

		return nil
	})
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

func (r *Repository) InitializeProducerProfile(userID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", userID, RoleProducer).Error; err != nil {
			return err
		}

		newProfile := producerProfile{
			UserId:   userID,
			Bio:      "",
			Websites: pq.StringArray{},
			Social:   SocialLinks{},
		}

		if err := tx.Create(&newProfile).Error; err != nil {
			return err
		}

		return nil
	})
}
