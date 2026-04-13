package exceptions

import "errors"

var ErrNotFound = errors.New("Entity Not Found")
var ErrEntityConflict = errors.New("Entity Conflict")
