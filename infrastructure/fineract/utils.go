package fineract

import "errors"

var (
	ErrInsufficientAccountBalance = errors.New("ErrInsufficientAccountBalance")
	ErrResourceNotFound           = errors.New("ErrResourceNotFound")
	ErrClientNotFound             = errors.New("ErrClientNotFound")
	errMap                        = map[string]error{
		"error.msg.savingsaccount.transaction.insufficient.account.balance": ErrInsufficientAccountBalance,
		"error.msg.resource.not.found":                                      ErrResourceNotFound,
		"error.msg.client.id.invalid":                                       ErrClientNotFound,
	}
)

func handleErrorCode(e string) error {
	if err, ok := errMap[e]; ok {
		return err
	}
	return errors.New(e)
}
