package exception

import "fmt"

type ErrorAppCode string

const (
	EntityNotFoundCode   ErrorAppCode = "ENTITY_NOT_FOUND"
	RoleAssignedCode     ErrorAppCode = "ROLE_ASSIGNED"
	EntityConflictCode   ErrorAppCode = "ENTITY_CONFLICT"
	ValidationFailedCode ErrorAppCode = "VALIDATION_FAILED"
	EmptyUpdateCode      ErrorAppCode = "EMPTY_UPDATE"
	InternalError        ErrorAppCode = "INTERNAL_ERROR"
)

type AppError struct {
	Message string       `json:"message"`
	Code    ErrorAppCode `json:"-"`
	Fields  []string     `json:"fields,omitempty"`
	Err     error        `json:"-"`
}

func (e *AppError) Type() string {
	return "[ APP ERROR ]"
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: \n%v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}
