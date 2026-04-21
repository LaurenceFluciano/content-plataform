package usecase

import (
	"errors"
	"fmt"

	"github.com/laurencefluciano/content-api/internal/app/user/command"
	"github.com/laurencefluciano/content-api/internal/domain/user"
	"github.com/laurencefluciano/content-api/internal/exception"
)

type UpdateProfileUseCase struct {
	Repo user.Repository
}

func (u UpdateProfileUseCase) Execute(authId string, cmd command.UpdateUserCommand) error {
	entity, err := u.Repo.FindByAuthId(authId)

	if err != nil {
		if errors.Is(err, exception.ErrRepoNotFound) {
			return &exception.AppError{
				Code:    exception.EntityNotFoundCode,
				Message: "Esse usuário não existe.",
				Err:     err,
			}
		}
		return err
	}

	cmd.Filter(entity.Roles())

	if cmd.IsEmpty() {
		return &exception.AppError{
			Code:    exception.EmptyUpdateCode,
			Message: "Nenhum campo para atualizar.",
		}
	}

	var errorFields []string

	// --- Nome ---
	if cmd.Name != nil {
		name, err := user.NewName(*cmd.Name)
		if err != nil {
			errorFields = append(errorFields, err.Error())
		} else {
			entity.ChangeName(name)
		}
	}

	// --- Avatar ---
	if cmd.AvatarUrl != nil {
		avatarUrl, err := user.NewAvatarUrl(*cmd.AvatarUrl)

		if err != nil {
			errorFields = append(errorFields, err.Error())
		} else {
			entity.ChangeAvatarUrl(avatarUrl)
		}
	}

	// --- Bio ---
	if cmd.Bio != nil {
		bio, err := user.NewBio(*cmd.Bio)
		if err != nil {
			errorFields = append(errorFields, err.Error())
		} else {
			entity.ChangeBio(bio)
		}
	}

	// --- Websites ---
	if cmd.Websites != nil {
		var validWebsites []user.Website
		has_invalid_url := false

		for i, raw := range *cmd.Websites {
			site, err := user.NewWebsite(raw)

			if err != nil {
				has_invalid_url = true
				errorFields = append(errorFields, fmt.Sprintf("website no índice %d é inválido: %s", i, err.Error()))
			}

			validWebsites = append(validWebsites, site)
		}

		if !has_invalid_url {
			err = entity.ChangeWebsites(validWebsites)
			if err != nil {
				errorFields = append(errorFields, err.Error())
			}
		}
	}

	// --- Social Links ---
	if cmd.Social != nil {
		entity.ChangeSocialLinks(*cmd.Social)
	}

	if len(errorFields) > 0 {
		return &exception.AppError{
			Code:    exception.ValidationFailedCode,
			Message: "Campos inválidos.",
			Fields:  errorFields,
		}
	}

	return u.Repo.Update(entity)
}
