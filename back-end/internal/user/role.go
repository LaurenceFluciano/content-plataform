package user

import "gorm.io/gorm"

type RoleType uint

const (
	RoleUnknown  RoleType = iota
	RoleReader   RoleType = 1
	RoleProducer RoleType = 2
)

type Role struct {
	RoleId RoleType `gorm:"primaryKey"`
	Name   string   `gorm:"uniqueIndex"`
}

func (u *User) HasRole(role RoleType) bool {
	for _, r := range u.Roles {
		if r.RoleId == role {
			return true
		}
	}
	return false
}

func CreateRoles(db *gorm.DB) error {
	roles := []Role{
		{RoleId: RoleReader, Name: "reader"},
		{RoleId: RoleProducer, Name: "producer"},
	}

	for _, role := range roles {
		err := db.FirstOrCreate(&role, Role{Name: role.Name}).Error
		if err != nil {
			return err
		}
	}

	return nil
}
