package dto

import (
	"time"

	"github.com/bekha-io/openbank/domain/types"
	"github.com/shopspring/decimal"
)

type CalculateAnnuityInstallmentsQuery struct {
	LoanAmount    types.Money
	InterestRate  decimal.Decimal // annual interest rate
	RepayStartsAt time.Time
	Duration      uint
}
