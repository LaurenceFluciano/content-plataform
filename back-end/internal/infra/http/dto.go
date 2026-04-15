package http

import "github.com/laurencefluciano/content-api/internal/domain/user"

type UpdateUserRequest struct {
	Name      *string           `json:"name"`
	AvatarUrl *string           `json:"avatar_url"`
	Bio       *string           `json:"bio"`
	Websites  *[]string         `json:"websites"`
	Social    *user.SocialLinks `json:"social"`
}

// Profile Response

type ProfileResponse struct {
	ID        string       `json:"id"`
	Name      *string      `json:"name,omitempty"`
	Roles     []string     `json:"roles"`
	AvatarUrl *string      `json:"avatar_url,omitempty"`
	Producer  *ProducerDTO `json:"producer,omitempty"`
}

type ProducerDTO struct {
	Bio      *string          `json:"bio,omitempty"`
	Websites []string         `json:"websites,omitempty"`
	Social   user.SocialLinks `json:"social,omitempty"`
}
