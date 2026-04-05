package user

type Role struct {
	RoleId uint   `gorm:"primaryKey"`
	Name   string `gorm:"uniqueIndex"`
}

func (Role) TableName() string {
	return "roles"
}
