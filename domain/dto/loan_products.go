package dto

import (
	"errors"
	"unicode/utf8"

	"github.com/bekha-io/openbank/domain/types"
	"github.com/shopspring/decimal"
)

type CreateLoanProductCommand struct {
	Name                string          // Product name
	MinDuration         uint            // Min loan duration (in months)
	MaxDuration         uint            // Max loan duration (in months)
	MinAmount           types.Money     // Min loan amount
	MaxAmount           types.Money     // Max loan amount
	InterestRate        decimal.Decimal // Interest rate
	LoanType            types.LoanType  // Annuity, etc.
	DailyOverduePenalty decimal.Decimal // Daily overdue penalty calculated on overdue unpaid installments
}


func (cmd CreateLoanProductCommand) Validate() error {
	var err error

	// len(name) < 3
	if utf8.RuneCountInString(cmd.Name) < 3 {
		err = errors.Join(err, errors.New("name length must be greater than 3"))
	}

	// min_duration > max_duration
	if cmd.MinDuration > cmd.MaxDuration {
		err = errors.Join(err, errors.New("min duration must be less than max duration"))
	}

	// min_amount > max_amount
	if cmd.MinAmount.Amount.GreaterThan(cmd.MaxAmount.Amount) {
		err = errors.Join(err, errors.New("max amount must be greater than min amount"))
	}

	// interest_rate < 0%
	if cmd.InterestRate.LessThan(decimal.NewFromInt(0)) {
		err = errors.Join(err, errors.New("interest rate must be greater than or equal to 0%"))
	}

	// daily_overdue_penalty < 0%
	if cmd.DailyOverduePenalty.LessThan(decimal.NewFromInt(0)) {
		err = errors.Join(err, errors.New("daily overdue penalty must be greater than or equal to 0%"))
	}

	// loan_type is invalid
	if !cmd.LoanType.IsValid() {
		err = errors.Join(err, errors.New("invalid loan type value"))
	}

	return err
}
