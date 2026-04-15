package user

import (
	"errors"
)

type RoleType uint

const (
	RoleUnknown  RoleType = iota
	RoleReader   RoleType = 1
	RoleProducer RoleType = 2
)

func (u *User) HasRole(role RoleType) bool {
	for _, r := range u.roles {
		if r == role {
			return true
		}
	}
	return false
}

func (u *User) AddRole(role RoleType) error {
	if u.HasRole(role) {
		return ErrRoleAlreadyAssigned
	}

	u.roles = append(u.roles, role)
	u.dirty[RolesField] = struct{}{}
	return nil
}

var ErrRoleAlreadyAssigned = errors.New("Role Already Assigned")
