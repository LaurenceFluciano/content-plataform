package user

import (
	"errors"

	"github.com/google/uuid"
)

type User struct {
	id      uuid.UUID
	authId  string
	name    Name
	status  Status
	roles   []RoleType
	profile *Profile

	dirty map[Field]struct{}
}

func (u *User) ID() uuid.UUID     { return u.id }
func (u *User) Name() string      { return u.name.Value() }
func (u *User) Status() Status    { return u.status }
func (u *User) Roles() []RoleType { return u.roles }
func (u *User) AuthId() string    { return u.authId }

func NewUser(authId string) (*User, error) {
	id, err := uuid.NewV7()

	if err != nil {
		return nil, errors.New("Erro ao gerar o id.")
	}

	return &User{
		id:     id,
		authId: authId,
		status: Active,
		roles:  []RoleType{RoleReader},
		dirty:  make(map[Field]struct{}),
	}, nil
}

func (u *User) ChangeName(name Name) {
	if name.Value() == u.Name() {
		return
	}

	u.name = name
	u.dirty[NameField] = struct{}{}
}

func (u *User) ChangeStatus(status Status) {
	u.status = status
	u.dirty[StatusField] = struct{}{}
}
