package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	AuthId    string
	Name      string
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
	Roles     []RoleType
	Profile   *Profile
}
