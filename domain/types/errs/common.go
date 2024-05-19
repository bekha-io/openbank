package errs

import "errors"

var (
	ErrNotAuthenticated = errors.New("Errors.Common.NotAuthenticated")
	ErrInternalDatabaseError = errors.New("Errors.Common.InternalDatabaseError")
)