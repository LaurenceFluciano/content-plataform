package dto

import "github.com/laurencefluciano/content-api/internal/domain/user"

type ProfileResponse struct {
	ID        string                   `json:"id"`
	Name      *string                  `json:"name,omitempty"`
	Roles     []string                 `json:"roles"`
	AvatarUrl *string                  `json:"avatar_url,omitempty"`
	Producer  *ProducerProfileResponse `json:"producer,omitempty"`
}

type ProducerProfileResponse struct {
	Bio      *string          `json:"bio,omitempty"`
	Websites []string         `json:"websites,omitempty"`
	Social   user.SocialLinks `json:"social,omitempty"`
}
