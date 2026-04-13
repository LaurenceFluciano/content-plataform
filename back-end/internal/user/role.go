package user

type RoleType uint

const (
	RoleUnknown  RoleType = iota
	RoleReader   RoleType = 1
	RoleProducer RoleType = 2
)

func (u *User) HasRole(role RoleType) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}
