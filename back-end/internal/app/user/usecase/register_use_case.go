package usecase

import (
	"errors"

	"github.com/laurencefluciano/content-api/internal/app/user/command"
	"github.com/laurencefluciano/content-api/internal/domain/user"
	"github.com/laurencefluciano/content-api/internal/exception"
)

type RegisterUserUseCase struct {
	Repo user.Repository
}

func (u RegisterUserUseCase) Execute(cmd command.RegisterUserCommand) error {
	entity, err := user.NewUser(cmd.AuthId)

	if err != nil {
		return &exception.AppError{
			Message: "Erro interno do servidor.",
			Code:    exception.InternalError,
			Err:     err,
		}
	}

	err = u.Repo.Create(entity)

	if err != nil {

		if errors.Is(err, exception.ErrRepoConflict) {
			return &exception.AppError{
				Code:    exception.EntityConflictCode,
				Message: "Usuário já existe.",
				Err:     err,
			}
		}

		return &exception.AppError{
			Code:    exception.InternalError,
			Message: "Erro interno do servidor.",
			Err:     err,
		}

	}

	return nil
}
