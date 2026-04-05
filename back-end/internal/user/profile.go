package user

import (
	"github.com/google/uuid"
)

type Profile struct {
	UserId     uuid.UUID `gorm:"primaryKey"`
	AvatartUrl string    `gorm:"type:varchar(255)"`
	User       *User
}

type Producer struct {
	UserId   uuid.UUID   `gorm:"primaryKey"`
	Bio      string      `gorm:"type:text"`
	Websites string      `gorm:"type:text"`
	Social   SocialLinks `gorm:"type:jsonb"`
	User     *User
}

func (Profile) TableName() string {
	return `profile`
}

func (Producer) TableName() string {
	return `producer_profile`
}
