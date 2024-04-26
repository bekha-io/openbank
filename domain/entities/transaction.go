package entities

import (
	"time"

	"github.com/bekha-io/vaultonomy/domain/types"
)

type Transaction struct {
	ID              types.TransactionID   `json:"id"`
	AccountID       types.AccountID       `json:"account_id"`
	TransactionType types.TransactionType `json:"transaction_type"`
	Amount          *types.Money          `json:"amount"`
	CreatedAt       time.Time             `json:"created_at"`
	Comment         string                `json:"comment"`
}

func NewTransaction(accountId types.AccountID, transactionType types.TransactionType, amount types.Money) *Transaction {
	return &Transaction{
		ID:              types.NewTransactionID(),
		AccountID:       accountId,
		TransactionType: transactionType,
		Amount:          &amount,
		CreatedAt:       time.Now().UTC(),
	}
}
