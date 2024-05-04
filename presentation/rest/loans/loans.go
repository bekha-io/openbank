package loans

import "github.com/bekha-io/openbank/domain/services"

type LoanController struct {
	LoanService services.ILoanService
}

func NewLoanController(loanSvc services.ILoanService) *LoanController {
	return &LoanController{
		LoanService: loanSvc,
	}
}
