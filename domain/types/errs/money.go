package errs

import "errors"

var (
	ErrMoneyIsZero = errors.New("Errors.Money.IsZero")
	ErrMoneyIsNegative = errors.New("Errors.Money.IsNegative")
)