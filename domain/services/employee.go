package services

import (
	"context"

	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/repository"
	"github.com/bekha-io/openbank/domain/types/errs"
)

type IEmployeeService interface {
	Authenticate(ctx context.Context, email, password string) (*entities.Employee, error)
	GetEmployeeByEmail(ctx context.Context, email string) (*entities.Employee, error)
	SearchEmployees(ctx context.Context, query string) ([]*entities.Employee, error)
	CreateEmployee(ctx context.Context, email, password, firstName, lastName, middleName string) error
}

var _ IEmployeeService = (*EmployeeService)(nil)

type EmployeeService struct {
	EmployeeRepo repository.IEmployeeRepository
}

func NewEmployeeService(employeeRepo repository.IEmployeeRepository) *EmployeeService {
	return &EmployeeService{
		EmployeeRepo: employeeRepo,
	}
}

// Authenticate implements IEmployeeService.
func (e *EmployeeService) Authenticate(ctx context.Context, email string, password string) (*entities.Employee, error) {
	emp, err := e.EmployeeRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !emp.IsPasswordCorrect(password) {
		return nil, errs.ErrEmployeeInvalidPassword
	}

	return emp, nil
}

// CreateEmployee implements IEmployeeService.
func (e *EmployeeService) CreateEmployee(ctx context.Context, email, password, firstName, lastName, middleName string) error {
	emp := entities.NewEmployee(email, password, firstName, lastName, middleName)
	return e.EmployeeRepo.Save(ctx, emp)
}

// GetEmployeeByEmail implements IEmployeeService.
func (e *EmployeeService) GetEmployeeByEmail(ctx context.Context, email string) (*entities.Employee, error) {
	return e.EmployeeRepo.GetByEmail(ctx, email)
}

// SearchEmployees implements IEmployeeService.
func (e *EmployeeService) SearchEmployees(ctx context.Context, query string) ([]*entities.Employee, error) {
	return e.EmployeeRepo.GetManyLike(ctx, query)
}
