package user

type Field string

const (
	NameField      Field = "name"
	StatusField    Field = "status"
	RolesField     Field = "roles"
	BioField       Field = "bio"
	WebsitesField  Field = "websites"
	SocialField    Field = "social"
	AvatarUrlField Field = "avatar_url"
)

func (u *User) MarkDirty(field Field) {
	if u.dirty == nil {
		u.dirty = make(map[Field]struct{})
	}
	u.dirty[field] = struct{}{}
}

func (u *User) GetDirtyFields() map[Field]struct{} {
	return u.dirty
}

func (u *User) ClearDirty() {
	u.dirty = make(map[Field]struct{})
}
