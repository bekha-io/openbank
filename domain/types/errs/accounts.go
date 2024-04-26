package errs

import "errors"

var (
	ErrAccountNotFound            = errors.New("Errors.Accounts.NotFound")
	ErrAccountDifferentCurrencies = errors.New("Errors.Accounts.DifferentCurrencies")
)
