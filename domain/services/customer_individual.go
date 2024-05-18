package services

import (
	"context"
	"errors"

	"github.com/bekha-io/openbank/domain/dto"
	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/repository"
	"github.com/bekha-io/openbank/domain/types"
	"github.com/bekha-io/openbank/domain/types/errs"
)

var _ IIndividualCustomerService = (*IndividualCustomerService)(nil)

type IndividualCustomerService struct {
	IndividualCustomerRepo repository.IIndividualCustomerRepository
	AccountsRepo           repository.IAccountRepository
}

func NewIndividualCustomerService(individualCustomerRepo repository.IIndividualCustomerRepository, accountsRepo repository.IAccountRepository) *IndividualCustomerService {
	return &IndividualCustomerService{
		IndividualCustomerRepo: individualCustomerRepo,
		AccountsRepo:           accountsRepo,
	}
}

// CreateCustomer implements IIndividualCustomerService.
func (i *IndividualCustomerService) CreateCustomer(ctx context.Context, in dto.CreateIndividualCustomerCommand) error {
	// If exists
	_, err := i.IndividualCustomerRepo.GetBy(ctx, "phone_number", in.PhoneNumber)
	if err == nil {
		return errors.New("Errors.Customers.PhoneNumberExists")
	}

	customer := entities.NewIndividualCustomer(in.PhoneNumber)
	err = i.IndividualCustomerRepo.Save(ctx, customer)
	return err
}

// GetCustomer implements IIndividualCustomerService.
func (i *IndividualCustomerService) GetCustomer(ctx context.Context, id types.CustomerID) (*entities.IndividualCustomer, error) {
	record, err := i.IndividualCustomerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Join(errs.ErrIndividualCustomerNotFound, err)
	}
	return record, err
}

// GetCustomerBy implements IIndividualCustomerService.
func (s *IndividualCustomerService) GetCustomerBy(ctx context.Context, key string, value interface{}) (*entities.IndividualCustomer, error) {
	record, err := s.IndividualCustomerRepo.GetBy(ctx, key, value)
	if err != nil {
		return nil, errors.Join(errs.ErrIndividualCustomerNotFound, err)
	}
	return record, err
}

// GetCustomerAccounts implements IIndividualCustomerService.
func (s *IndividualCustomerService) GetCustomerAccounts(ctx context.Context, customer entities.Customer) ([]*entities.Account, error) {
	accounts, err := s.AccountsRepo.GetManyBy(ctx, repository.Filter{Key: "customer_id", EqualTo: customer.Id().String()})
	if err != nil {
		return nil, errors.Join(errs.ErrIndividualCustomerNotFound, err)
	}
	return accounts, nil
}

// GetCustomersLike implements IIndividualCustomerService.
func (i *IndividualCustomerService) SearchCustomersByPhoneNumber(ctx context.Context, phoneNumber string) ([]*entities.IndividualCustomer, error) {
	return i.IndividualCustomerRepo.GetManyPhoneNumberLike(ctx, phoneNumber)
}
