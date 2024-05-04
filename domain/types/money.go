package types

import (
	"fmt"

	"github.com/bekha-io/openbank/domain/types/errs"
	"github.com/shopspring/decimal"
)

type PaymentType string

const (
	PaymentTypeCash     PaymentType = "cash"
	PaymentTypeCashless PaymentType = "cashless"
)

type Money struct {
	Amount   decimal.Decimal `json:"amount"`
	Currency Currency        `json:"currency"`
}

func NewMoney(amount decimal.Decimal, currency Currency) *Money {
	return &Money{
		Amount:   amount,
		Currency: currency,
	}
}

func (m *Money) String() string {
	return fmt.Sprintf("%v %v", m.Amount.StringFixedBank(2), m.Currency)
}

func (m *Money) Add(mm Money) error {
	m.Amount = m.Amount.Add(mm.Amount)
	return nil
}

func (m *Money) Sub(mm Money) error {
	m.Amount = m.Amount.Sub(mm.Amount)
	return nil
}

func (m Money) Validate() error {
	if !m.Amount.IsPositive() {
		return errs.ErrMoneyIsNegative
	}

	if m.Amount.IsZero() {
		return errs.ErrMoneyIsZero
	}

	return nil
}
