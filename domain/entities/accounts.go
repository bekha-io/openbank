package entities

import (
	"time"

	"github.com/bekha-io/openbank/domain/types"
	"github.com/shopspring/decimal"
)

type Account struct {
	ID         types.AccountID  `json:"id"`
	CustomerID types.CustomerID `json:"customer_id"`
	Balance    *types.Money     `json:"balance"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
}

func NewAccount(customerId types.CustomerID, currency types.Currency) *Account {
	return &Account{
		ID:         types.NewAccountID(),
		CustomerID: customerId,
		Balance:    types.NewMoney(decimal.NewFromInt(0), currency),
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}
