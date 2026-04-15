package usecase

import (
	"path/filepath"
	"strings"

	"github.com/laurencefluciano/content-api/internal/app/exception"
	"github.com/laurencefluciano/content-api/internal/domain/user"
)

type UpdateProfileUseCase struct {
	Repo user.Repository
}

type UpdateUserParams struct {
	Name      *string
	AvatarUrl *string
	Bio       *string
	Websites  *[]string
	Social    *user.SocialLinks
}

func (u UpdateProfileUseCase) Execute(authId string, params UpdateUserParams) error {
	profile, err := u.Repo.FindByAuthId(authId)
	if err != nil {
		return err
	}

	allowedFields := user.GetEditableFieldsMap(profile.Roles())
	var fieldErrors []string

	if params.Name != nil && allowedFields["Name"] {
		newName, err := user.NewName(*params.Name)
		if err != nil {
			fieldErrors = append(fieldErrors, err.Error())
		} else {
			existing, err := u.Repo.FindByName(newName.Value())
			if err == nil && existing.AuthId() != authId {
				fieldErrors = append(fieldErrors, "Este nome de usuário já está sendo utilizado.")
			} else {
				profile.SetName(newName.Value())
			}
		}
	}

	if params.AvatarUrl != nil && allowedFields["AvatarUrl"] {
		if !isValidStoragePath(*params.AvatarUrl) {
			fieldErrors = append(fieldErrors, "URL da imagem do Avatar é inválida.")
		} else {
			profile.SetAvatarUrl(*params.AvatarUrl)
		}
	}

	// Atenção isso deve virar um vo
	if params.Bio != nil && allowedFields["Bio"] {
		if len(*params.Bio) < 10 {
			fieldErrors = append(fieldErrors, "A bio deve ter pelo menos 10 caracteres.")
		} else {
			profile.SetBio(*params.Bio)
		}
	}

	if params.Websites != nil && allowedFields["Websites"] {
		for _, website := range *params.Websites {
			err = profile.AddWebsite(website)

			if err != nil {
				fieldErrors = append(fieldErrors, err.Error())
				break
			}

		}
	}

	if params.Social != nil && allowedFields["Social"] {
		profile.ChangeSocialLinks(*params.Social)
	}

	if len(fieldErrors) > 0 {
		return &exception.AppError{
			Code:   exception.ValidationFailedCode,
			Fields: fieldErrors,
		}
	}

	if !profile.IsDirty() {
		return &exception.AppError{
			Code:   exception.EmptyUpdateCode,
			Fields: fieldErrors,
		}
	}

	return u.Repo.Save(profile)
}

func isValidStoragePath(path string) bool {
	if path == "" {
		return false
	}

	if !strings.HasPrefix(path, "avatars/") {
		return false
	}

	ext := strings.ToLower(filepath.Ext(path))
	validExtensions := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
	}

	return validExtensions[ext]
}
