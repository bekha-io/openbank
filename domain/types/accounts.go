package types

import (
	"log/slog"

	"github.com/jacoelho/banking/iban"
)

type AccountID string

func NewAccountID() AccountID {
	id, err := iban.Generate("AE")
	if err != nil {
		slog.Error(err.Error())
	}
	return AccountID(id)
}
