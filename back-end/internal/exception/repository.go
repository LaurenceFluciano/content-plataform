package exception

import (
	"fmt"
)

type ErrorRepositoryCode string

const (
	NotFound       ErrorRepositoryCode = "NOT_FOUND"
	Conflict       ErrorRepositoryCode = "CONFLICT"
	InvalidId      ErrorRepositoryCode = "INVALID_ID"
	FailConnection ErrorRepositoryCode = "FAIL_CONNECTION"
)

type ErrorRepository struct {
	Code    ErrorRepositoryCode `json:"-"`
	Message string              `json:"message"`
	Err     error               `json:"-"`
}

func (e *ErrorRepository) Type() string {
	return "[ REPOSITORY ERROR ]"
}

func (e *ErrorRepository) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: \n%v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *ErrorRepository) Unwrap() error {
	return e.Err
}

var (
	ErrRepoNotFound       = &ErrorRepository{Code: NotFound}
	ErrRepoConflict       = &ErrorRepository{Code: Conflict}
	ErrRepoInvalidId      = &ErrorRepository{Code: InvalidId}
	ErrRepoFailConnection = &ErrorRepository{Code: FailConnection}
)

func (e *ErrorRepository) Is(target error) bool {
	t, ok := target.(*ErrorRepository)
	if !ok {
		return false
	}
	return e.Code == t.Code
}
