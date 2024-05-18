package employees

import "github.com/bekha-io/openbank/domain/services"

type EmployeeController struct {
	EmployeeService services.IEmployeeService
}

func NewEmployeeController(empSvc services.IEmployeeService) *EmployeeController {
	return &EmployeeController{
		EmployeeService: empSvc,
	}
}
