package usecase

import (
	"errors"

	"github.com/laurencefluciano/content-api/internal/app/user/command"
	"github.com/laurencefluciano/content-api/internal/domain/user"
	"github.com/laurencefluciano/content-api/internal/exception"
)

type BecomeProducerUserCase struct {
	Repo user.Repository
}

func (u *BecomeProducerUserCase) Execute(cmd command.BecomeProducerCommand) error {
	entity, err := u.Repo.FindByAuthId(cmd.AuthId)

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

	err = entity.AddRole(user.RoleProducer)

	if err != nil {
		if errors.Is(err, user.ErrRoleAlreadyAssigned) {
			return &exception.AppError{
				Code:    exception.RoleAssignedCode,
				Message: "Esse usuário já tem esse papel.",
				Err:     err,
			}
		}
		return err
	}

	return u.Repo.Update(entity)
}
