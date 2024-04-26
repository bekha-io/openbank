package entities

import "github.com/bekha-io/vaultonomy/domain/types"

type Loan interface {
	TotalAmount() types.Money
	LoanProduct() LoanProduct
	Duration() uint // Loan duration in days
	Installments() []types.Installment
}
