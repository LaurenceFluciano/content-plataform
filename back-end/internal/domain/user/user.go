package user

import (
	"errors"

	"github.com/google/uuid"
)

type Field string

const (
	NameField      Field = "name"
	StatusField    Field = "status"
	RolesField     Field = "roles"
	BioField       Field = "bio"
	WebsitesField  Field = "website"
	SocialField    Field = "social"
	AvatarUrlField Field = "avatar_url"
)

type User struct {
	id      uuid.UUID
	authId  string
	name    Name
	status  Status
	roles   []RoleType
	profile *Profile

	dirty map[Field]struct{}
	isNew bool
}

func NewUser(authId string) (*User, error) {
	id, err := uuid.NewV7()

	if err != nil {
		return nil, errors.New("Erro ao gerar o id.")
	}

	return &User{
		id:     id,
		authId: authId,
		isNew:  true,
		dirty:  make(map[Field]struct{}),
	}, nil
}

func (u *User) SetName(name string) error {
	nameVo, err := NewName(name)

	if err != nil {
		return err
	}

	if nameVo.Value() == u.Name() {
		return nil
	}

	u.name = nameVo
	u.dirty[NameField] = struct{}{}

	return nil
}

func (u *User) SetStatus(status Status) {
	u.status = status
	u.dirty[StatusField] = struct{}{}
}

func (u *User) Name() string {
	return u.name.Value()
}

func (u *User) Status() Status {
	return u.status
}

func (u *User) Roles() []RoleType {
	return u.roles
}

func (u *User) Profile() *Profile {
	return u.profile
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) AuthId() string {
	return u.authId
}

func (u *User) IsNew() bool {
	return u.isNew
}

func (u *User) GetDirtyFields() map[Field]struct{} {
	return u.dirty
}

func (u *User) IsDirty() bool {
	return len(u.dirty) > 0
}

func (u *User) SetPersisted() {
	u.isNew = false
	u.dirty = make(map[Field]struct{})
}

func RestoreUser(id uuid.UUID, authId string, name string, status Status) *User {
	return &User{
		id:     id,
		authId: authId,
		name:   Name{value: name},
		status: status,
		dirty:  make(map[Field]struct{}),
		isNew:  false,
	}
}

func (u *User) RestoreProfile(avatarUrl string) {
	if u.profile == nil {
		u.profile = &Profile{}
	}
	u.profile.avatarUrl = avatarUrl
}

func (u *User) RestoreProducer(bio string, websites []string, social SocialLinks) {
	if u.profile == nil {
		u.profile = &Profile{}
	}
	u.profile.bio = bio
	u.profile.websites = websites
	u.profile.social = social
}
