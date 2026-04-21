package repository

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/laurencefluciano/content-api/internal/exception"
	"gorm.io/gorm"
)

func (r *gormUserRepository) handleDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &exception.ErrorRepository{
			Code:    exception.NotFound,
			Message: "O registro solicitado não foi encontrado no banco de dados.",
			Err:     err,
		}
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return &exception.ErrorRepository{
			Code:    exception.Conflict,
			Message: fmt.Sprintf("Conflito de unicidade na restrição: %s", pgErr.ConstraintName),
			Err:     err,
		}
	}

	return &exception.ErrorRepository{
		Code:    exception.FailConnection,
		Message: "Falha na comunicação com o banco de dados.",
		Err:     err,
	}
}
