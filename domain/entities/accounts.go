package entities

import (
	"time"

	"github.com/bekha-io/openbank/domain/types"
)

type Account struct {
	ID         uint         `json:"id"`
	AccountNo  string       `json:"account_no"`
	CustomerID uint         `json:"customer_id"`
	Balance    *types.Money `json:"balance"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}
