package accounts

import "github.com/bekha-io/vaultonomy/domain/services"

type AccountsController struct {
	AccountsService services.IAccountService
}

func NewAccountsController(accountsSvc services.IAccountService) *AccountsController {
	return &AccountsController{
		AccountsService: accountsSvc,
	}
}
