package customers

import "github.com/bekha-io/openbank/domain/services"

type CustomerController struct {
	IndividualCustomerService services.ICustomerService
}

func NewCustomerController(individualCustomerService services.ICustomerService) *CustomerController {
	return &CustomerController{
		IndividualCustomerService: individualCustomerService,
	}
}
