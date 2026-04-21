package dto

import "github.com/laurencefluciano/content-api/internal/domain/user"

type UpdateUserRequest struct {
	Name      *string           `json:"name"`
	AvatarUrl *string           `json:"avatar_url"`
	Bio       *string           `json:"bio"`
	Websites  *[]string         `json:"websites"`
	Social    *user.SocialLinks `json:"social"`
}
