package customers

import "github.com/bekha-io/vaultonomy/domain/services"

type CustomerController struct {
	IndividualCustomerService services.IIndividualCustomerService
}

func NewCustomerController(individualCustomerService services.IIndividualCustomerService) *CustomerController {
	return &CustomerController{
		IndividualCustomerService: individualCustomerService,
	}
}
