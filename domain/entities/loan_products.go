package entities

import (
	"errors"
	"time"

	"github.com/bekha-io/openbank/domain/types"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// LoanProduct represents loan product with customizable terms and conditions
type LoanProduct struct {
	ID                  types.LoanProductID `json:"id"`                    // Product ID (custom identifier)
	Name                string              `json:"name"`                  // Product name
	MinDuration         uint                `json:"min_duration"`          // Min loan duration (in months)
	MaxDuration         uint                `json:"max_duration"`          // Max loan duration (in months)
	MinAmount           types.Money         `json:"min_amount"`            // Min loan amount
	MaxAmount           types.Money         `json:"max_amount"`            // Max loan amount
	InterestRate        decimal.Decimal     `json:"interest_rate"`         // Interest rate
	LoanType            types.LoanType      `json:"loan_type"`             // Annuity, etc.
	DailyOverduePenalty decimal.Decimal     `json:"daily_overdue_penalty"` // Daily overdue penalty calculated on overdue unpaid installments
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
}

func NewLoanProduct(name string, minAmount, maxAmount decimal.Decimal, currency types.Currency, interestRate decimal.Decimal, loanType types.LoanType, dailyOverduePenalty decimal.Decimal, minDuration uint, maxDuration uint) (*LoanProduct, error) {
	if minAmount.GreaterThan(maxAmount) {
		return nil, errors.New("max amount must be greater than min amount")
	}

	if interestRate.LessThan(decimal.NewFromInt(0)) {
		return nil, errors.New("interest rate must be greater than or equal to 0%")
	}

	return &LoanProduct{
		ID:                  types.LoanProductID(uuid.New()),
		Name:                name,
		MinAmount:           *types.NewMoney(minAmount, currency),
		MaxAmount:           *types.NewMoney(maxAmount, currency),
		InterestRate:        interestRate,
		LoanType:            loanType,
		MinDuration: minDuration,
		MaxDuration: maxDuration,
		DailyOverduePenalty: dailyOverduePenalty,
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
	}, nil
}
