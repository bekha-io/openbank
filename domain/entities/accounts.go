package entities

import (
	"github.com/bekha-io/openbank/domain/types"
	"github.com/shopspring/decimal"
)

type Account struct {
	ID         types.AccountID  `json:"id"`
	CustomerID types.CustomerID `json:"customer_id"`
	Balance    *types.Money     `json:"balance"`

	TrustedCustomers []*Customer `json:"-"` // Customers that are allowed to access this account and perform operations
}

func NewAccount(customerId types.CustomerID, currency types.Currency) *Account {
	return &Account{
		ID:         types.NewAccountID(),
		CustomerID: customerId,
		Balance:    types.NewMoney(decimal.NewFromInt(0), currency),
	}
}
