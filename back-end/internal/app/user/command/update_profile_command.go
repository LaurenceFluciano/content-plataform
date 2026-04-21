package command

import (
	"reflect"

	"github.com/laurencefluciano/content-api/internal/domain/user"
)

type UpdateUserCommand struct {
	Name      *string           `field:"name"`
	AvatarUrl *string           `field:"avatar_url"`
	Bio       *string           `field:"bio"`
	Websites  *[]string         `field:"websites"`
	Social    *user.SocialLinks `field:"social"`
}

func (c *UpdateUserCommand) Filter(roles []user.RoleType) {
	v := reflect.ValueOf(c).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldTag := t.Field(i).Tag.Get("field")
		if fieldTag == "" {
			continue
		}

		if !user.CanChangeField(roles, user.Field(fieldTag)) {
			v.Field(i).Set(reflect.Zero(v.Field(i).Type()))
		}
	}
}

func (c *UpdateUserCommand) IsEmpty() bool {
	return c.Name == nil &&
		c.AvatarUrl == nil &&
		c.Bio == nil &&
		c.Websites == nil &&
		c.Social == nil
}
