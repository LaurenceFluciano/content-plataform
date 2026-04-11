package user

type CreateUserRequest struct {
	Name *string `json:"name"`
}

// Profile Response

type ProfileResponse struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	AvatarUrl  string       `json:"avatar_url,omitempty"`
	IsProducer bool         `json:"is_producer"`
	Producer   *ProducerDTO `json:"producer,omitempty"`
}

type ProducerDTO struct {
	Bio      string      `json:"bio"`
	Websites string      `json:"websites"`
	Social   SocialLinks `json:"social"`
}
