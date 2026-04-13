package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/laurencefluciano/content-api/internal/user"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserTable struct {
	ID        uuid.UUID   `gorm:"primaryKey;column:id;default:gen_random_uuid();->"`
	AuthId    string      `gorm:"type:varchar(256);uniqueIndex"`
	Name      string      `gorm:"type:varchar(100);uniqueIndex"`
	Status    user.Status `gorm:"type:varchar(20);default:pending"`
	CreatedAt time.Time   `gorm:"autoCreateTime"`
	UpdatedAt time.Time   `gorm:"autoUpdateTime"`
	Roles     []RoleTable `gorm:"many2many:user_roles;joinForeignKey:user_id;JoinReferences:role_id"`

	UserProfile     *UserProfileTable     `gorm:"foreignKey:UserId"`
	ProducerProfile *ProducerProfileTable `gorm:"foreignKey:UserId"`
}

type RoleTable struct {
	RoleId user.RoleType `gorm:"primaryKey"`
	Name   string        `gorm:"uniqueIndex"`
}

type UserProfileTable struct {
	ID        uuid.UUID  `gorm:"primaryKey;column:user_id"`
	AvatarUrl string     `gorm:"type:varchar(255);column:avatar_url"`
	User      *UserTable `gorm:"foreignKey:UserId"`
}

type ProducerProfileTable struct {
	ID       uuid.UUID        `gorm:"primaryKey;column:user_id"`
	Bio      string           `gorm:"type:text"`
	Websites pq.StringArray   `gorm:"type:text[]"`
	Social   user.SocialLinks `gorm:"type:jsonb"`
	User     *UserTable       `gorm:"foreignKey:UserId"`
}

func GetModels() []interface{} {
	return []interface{}{
		&UserTable{},
		&RoleTable{},
		&UserProfileTable{},
		&ProducerProfileTable{},
	}
}

func CreateRoles(db *gorm.DB) error {
	roles := []RoleTable{
		{RoleId: user.RoleReader, Name: "reader"},
		{RoleId: user.RoleProducer, Name: "producer"},
	}

	for _, role := range roles {
		err := db.FirstOrCreate(&role, RoleTable{Name: role.Name}).Error
		if err != nil {
			return err
		}
	}

	return nil
}
