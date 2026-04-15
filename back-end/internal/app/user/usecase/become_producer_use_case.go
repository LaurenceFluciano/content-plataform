package usecase

import (
	"errors"

	"github.com/laurencefluciano/content-api/internal/app/exception"
	"github.com/laurencefluciano/content-api/internal/domain/user"
)

type BecomeProducerUserCase struct {
	Repo user.Repository
}

func (u *BecomeProducerUserCase) Execute(authId string) error {
	entity, err := u.Repo.FindByAuthId(authId)

	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return &exception.AppError{
				Code:    exception.EntityNotFoundCode,
				Message: "Esse usuário não existe.",
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
			}
		}
		return err
	}

	return u.Repo.Save(entity)
}
