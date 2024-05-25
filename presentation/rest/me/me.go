package me

import "github.com/bekha-io/openbank/domain/services"

type Controller struct {
	AccountsService  services.IAccountService
	CustomersService services.IIndividualCustomerService
}

func NewController(accountsSvc services.IAccountService, customersSvc services.IIndividualCustomerService) *Controller {
	return &Controller{
		AccountsService:  accountsSvc,
		CustomersService: customersSvc,
	}
}
