package types

import (
	"encoding/json"

	"github.com/google/uuid"
)

type TransactionID uuid.UUID

func (t TransactionID) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t TransactionID) String() string {
	return uuid.UUID(t).String()
}

func (t TransactionID) UUID() uuid.UUID {
	return uuid.UUID(t)
}

func NewTransactionID() TransactionID {
	return TransactionID(uuid.New())
}

type TransactionType string

const (
	DepositTransactionType  TransactionType = "deposit"
	WithdrawTransactionType TransactionType = "withdraw"
)
