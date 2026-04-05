package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"primaryKey;column:id"`
	AuthId    string    `gorm:"type:varchar(256);uniqueIndex"`
	Name      string    `gorm:"type:varchar(100);"`
	Status    Status    `gorm:"type:varchar(20);default:pending"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Roles     []Role    `gorm:"many2many:user_roles;"`
}

func (User) TableName() string {
	return "user"
}
