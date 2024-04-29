package loans

import "github.com/bekha-io/vaultonomy/domain/services"

type LoanController struct {
	LoanService services.ILoanService
}

func NewLoanController(loanSvc services.ILoanService) *LoanController {
	return &LoanController{
		LoanService: loanSvc,
	}
}
