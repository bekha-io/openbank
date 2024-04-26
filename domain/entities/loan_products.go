package entities

import "github.com/bekha-io/vaultonomy/domain/types"

// LoanProductEligibilityChecker represents a callback function that checks whether `c` is eligible for loan product `lp`
type LoanProductEligibilityChecker func(lp LoanProduct, c IndividualCustomer) bool

// LoanProduct represents loan product with customizable terms and conditions
type LoanProduct struct {
	Name        string      // Product name
	MinDuration uint        // Min loan duration (in days)
	MaxDuration uint        // Max loan duration (in days)
	MinAmount   types.Money // Min loan amount
	MaxAmount   types.Money // Max loan amount
}
