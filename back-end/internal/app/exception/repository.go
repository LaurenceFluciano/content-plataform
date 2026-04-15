package exception

import "errors"

var ErrNotFound = errors.New("Entity Not Found")
var ErrEntityConflict = errors.New("Entity Conflict")
var ErrInvalidEntityID = errors.New("Invalid ID")
var ErrInternal = errors.New("internal server error")
