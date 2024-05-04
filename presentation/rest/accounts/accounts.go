package accounts

import "github.com/bekha-io/openbank/domain/services"

type AccountsController struct {
	AccountsService services.IAccountService
}

func NewAccountsController(accountsSvc services.IAccountService) *AccountsController {
	return &AccountsController{
		AccountsService: accountsSvc,
	}
}
