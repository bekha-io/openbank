package me

import "github.com/bekha-io/openbank/domain/services"

type Controller struct {
	AccountsService  services.IAccountService
	CustomersService services.ICustomerService
}

func NewController(accountsSvc services.IAccountService, customersSvc services.ICustomerService) *Controller {
	return &Controller{
		AccountsService:  accountsSvc,
		CustomersService: customersSvc,
	}
}
