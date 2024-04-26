package types

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestCurrencyRate_ConvertTo(t *testing.T) {

	cr := CurrencyRate{
		Currency: "USD",
		Rates: map[Currency]decimal.Decimal{
			"TJS": decimal.NewFromFloat(10),
			"RUB": decimal.NewFromFloat(91.91),
		},
	}

	m1 := &Money{
		Amount: decimal.NewFromInt(100),
		Currency: "TJS",
	}
	result, _ := cr.Convert(m1)
	expected := &Money{Currency: "USD", Amount: decimal.NewFromInt(10)}
	require.Equal(t, expected.Amount.StringFixed(2), result.Amount.StringFixed(2))

	m2 := &Money{
		Amount: decimal.NewFromFloat(1000),
		Currency: "RUB",
	}
	result, _ = cr.Convert(m2)
	expected = &Money{Currency: "USD", Amount: decimal.NewFromFloat(10.88)}
	require.Equal(t, expected.Amount.StringFixed(2), result.Amount.StringFixed(2))
}
