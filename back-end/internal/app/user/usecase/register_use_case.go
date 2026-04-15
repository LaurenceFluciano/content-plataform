package usecase

import (
	"errors"

	"github.com/laurencefluciano/content-api/internal/app/exception"
	"github.com/laurencefluciano/content-api/internal/domain/user"
)

type RegisterUserUseCase struct {
	Repo user.Repository
}

func (u RegisterUserUseCase) Execute(authId string) error {
	userEntity, err := u.Repo.FindByAuthId(authId)

	if err == nil {
		return &exception.AppError{
			Code:    exception.EntityConflictCode,
			Message: "Usuário já existe.",
		}
	}

	if errors.Is(err, exception.ErrNotFound) {
		userEntity, err = user.NewUser(authId)
		if err != nil {
			return err
		}

		userEntity.SetStatus(user.Active)

		return u.Repo.Save(userEntity)
	}

	userEntity.AddRole(user.RoleReader)

	return err
}
