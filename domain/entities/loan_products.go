package entities

import (
	"github.com/bekha-io/openbank/domain/types"
	"github.com/shopspring/decimal"
)

// LoanProductEligibilityChecker represents a callback function that checks whether `c` is eligible for loan product `lp`
type LoanProductEligibilityChecker func(lp LoanProduct, c IndividualCustomer) error

// LoanProduct represents loan product with customizable terms and conditions
type LoanProduct struct {
	ID           types.LoanProductID // Product ID (custom identifier)
	Name         string              // Product name
	MinDuration  uint                // Min loan duration (in months)
	MaxDuration  uint                // Max loan duration (in months)
	MinAmount    types.Money         // Min loan amount
	MaxAmount    types.Money         // Max loan amount
	InterestRate decimal.Decimal
}
