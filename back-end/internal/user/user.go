package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"primaryKey;column:id;default:gen_random_uuid();->"`
	AuthId    string    `gorm:"type:varchar(256);uniqueIndex"`
	Name      string    `gorm:"type:varchar(100);uniqueIndex"`
	Status    Status    `gorm:"type:varchar(20);default:pending"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Roles     []Role    `gorm:"many2many:user_roles;joinForeignKey:user_id;JoinReferences:role_id"`

	UserProfile     *userProfile     `gorm:"foreignKey:UserId"`
	ProducerProfile *producerProfile `gorm:"foreignKey:UserId"`
}

func GetModels() []interface{} {
	return []interface{}{
		&User{},
		&Role{},
		&userProfile{},
		&producerProfile{},
	}
}
