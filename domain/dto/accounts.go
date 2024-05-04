package dto

import (
	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/types"
)

type CreateAccountCommand struct {
	CustomerID types.CustomerID
	Currency   types.Currency
}

type DepositCommand struct {
	Account *entities.Account
	Money   types.Money
	Comment string
}

type WithdrawCommand struct {
	Account *entities.Account
	Money   types.Money
	Comment string
}
