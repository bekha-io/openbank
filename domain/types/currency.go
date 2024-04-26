package types

import (
	"github.com/bekha-io/vaultonomy/domain/types/errs"
	"github.com/shopspring/decimal"
)

type Currency string

type CurrencyRate struct {
	Currency Currency
	Rates    map[Currency]decimal.Decimal  // each rate represents an amount of a given currency to 1 CurrencyRate.Currency
}

// Convert converts m.Currency to CurrencyRate.Currency based on CurrencyRate.Rates
func (c CurrencyRate) Convert(m *Money) (*Money, error) {
	// If currencies are the same
	if m.Currency == c.Currency {
		return &Money{
			Currency: m.Currency, Amount: decimal.NewFromInt(1),
		}, nil
	}

	rate, ok := c.Rates[m.Currency]
	if !ok {
		return nil, errs.ErrCurrencyRateNotFound
	}

	return &Money{
		Currency: c.Currency,
		Amount:   m.Amount.Div(rate),
	}, nil
}
