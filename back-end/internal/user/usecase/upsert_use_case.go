package usecase

import (
	"errors"
	"strings"

	"github.com/laurencefluciano/content-api/internal/user"
)

type UseCase struct {
	repository user.Repository
}

type User struct {
	Name      *string
	AvatarUrl *string
	Bio       *string
	Websites  []string
	Social    *user.SocialLinks
}

func UpsertUserProfileUseCase(params User) error {

	if params.Name != nil {
		if len(*params.Name) < 2 {
			return errors.New("Nome deve ter no minimo 2 caracteres.")
		}
		*params.Name = strings.ToLower(*params.Name)
	}

	return nil
}
