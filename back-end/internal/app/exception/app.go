package exception

import "fmt"

type ErrorAppCode string

const (
	EntityNotFoundCode   string = "ENTITY_NOT_FOUND"
	RoleAssignedCode     string = "ROLE_ASSIGNED"
	EntityConflictCode   string = "ENTITY_CONFLICT"
	ValidationFailedCode string = "VALIDATION_FAILED"
	EmptyUpdateCode      string = "EMPTY_UPDATE"
)

type AppError struct {
	Message string
	Err     error
	Code    string
	Fields  []string
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}
