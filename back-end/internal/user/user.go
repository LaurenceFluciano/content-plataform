package user

import (
	"time"
)

type User struct {
	Id        string
	AuthId    string
	Name      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserAggregate struct {
	User     User
	Profile  *Profile
	Producer *Producer
}
