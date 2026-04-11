package user

import (
	"github.com/google/uuid"
)

type Profile struct {
	User            *User
	UserProfile     *userProfile
	ProducerProfile *producerProfile
}

type userProfile struct {
	UserId    uuid.UUID `gorm:"primaryKey;column:user_id"`
	AvatarUrl string    `gorm:"type:varchar(255);column:avatar_url"`
	User      *User     `gorm:"foreignKey:UserId"`
}

type producerProfile struct {
	UserId   uuid.UUID   `gorm:"primaryKey;column:user_id"`
	Bio      string      `gorm:"type:text"`
	Websites string      `gorm:"type:text"`
	Social   SocialLinks `gorm:"type:jsonb"`
	User     *User       `gorm:"foreignKey:UserId"`
}
