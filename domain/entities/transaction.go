package entities

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransactionStatus string

const (
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

type Transaction struct {
	ID            uint              `json:"id"`
	FromAccountId uint              `json:"from_account_id"`
	ToAccountId   uint              `json:"to_account_id"`
	Status        TransactionStatus `json:"status"`
	StatusReason  string            `json:"status_reason"`
	Comment       string            `json:"comment"`
	Amount        decimal.Decimal   `json:"amount"`
	CreatedAt     time.Time         `json:"created_at"`
}
