package services

import (
	"context"
	"errors"
	"time"

	"github.com/bekha-io/openbank/domain/dto"
	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/repository"
	"github.com/bekha-io/openbank/domain/types/errs"
	"github.com/bekha-io/openbank/infrastructure/fineract"
)

var _ ICustomerService = (*IndividualCustomerService)(nil)

type IndividualCustomerService struct {
	Core            fineract.FineractClientI
	BeneficiaryRepo repository.IBenificiaryRepository
}

func NewIndividualCustomerService(core fineract.FineractClientI, beneficiaryRepo repository.IBenificiaryRepository) *IndividualCustomerService {
	return &IndividualCustomerService{
		Core:            core,
		BeneficiaryRepo: beneficiaryRepo,
	}
}

// CreateCustomer implements IIndividualCustomerService.
func (i *IndividualCustomerService) CreateCustomer(ctx context.Context, in dto.CreateIndividualCustomerCommand) error {
	// If exists
	_, err := i.Core.GetClientByExternalId(ctx, in.PhoneNumber)
	if err == nil {
		return errors.New("Errors.Customers.PhoneNumberExists")
	}

	// Create
	_, err = i.Core.CreateClient(ctx, fineract.CreateClientIn{
		ActivationDate:   time.Now().Format("02.01.2006"),
		DateFormat:       "dd.MM.yyyy",
		Locale:           "ru",
		ExternalId:       in.PhoneNumber,
		MobileNo:         in.PhoneNumber,
		OfficeId:         1,
		SavingsProductId: 1,
		FirstName:        "Unverified",
		LastName:         "User",
		Active:           true,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetCustomer implements IIndividualCustomerService.
func (i *IndividualCustomerService) GetCustomer(ctx context.Context, id uint) (*entities.Customer, error) {
	r, err := i.Core.GetClientById(ctx, id)
	if err != nil {
		return nil, errors.Join(errs.ErrIndividualCustomerNotFound, err)
	}
	return &entities.Customer{
		ID:          r.Id,
		PhoneNumber: r.ExternalId,
		FirstName:   r.Firstname,
		LastName:    r.Lastname,
	}, err
}

// GetCustomerAccounts implements IIndividualCustomerService.
func (s *IndividualCustomerService) GetCustomerAccounts(ctx context.Context, customer entities.Customer) ([]*entities.Account, error) {
	fAccounts, err := s.Core.GetClientsAccounts(ctx, customer.ID)
	if err != nil {
		return nil, errors.Join(errs.ErrIndividualCustomerNotFound, err)
	}

	var accounts = []*entities.Account{}

	// We need to extract balances as well (EXPENSIVE OPERATION)
	for _, account := range fAccounts.SavingsAccounts {
		// Assuming we have a separate service for balance retrieval
		detail, err := s.Core.GetSavingsAccountById(ctx, account.ID)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, detail.ToEntity())
	}

	return accounts, nil
}

func (s *IndividualCustomerService) GetCustomerByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.Customer, error) {
	r, err := s.Core.GetClientByExternalId(ctx, phoneNumber)
	if err != nil {
		return nil, errors.Join(errs.ErrIndividualCustomerNotFound, err)
	}
	return &entities.Customer{
		ID:          r.Id,
		PhoneNumber: r.ExternalId,
		FirstName:   r.Firstname,
		LastName:    r.Lastname,
	}, err
}

// GetBenificiaryByID implements ICustomerService.
func (i *IndividualCustomerService) GetBeneficiaryByID(ctx context.Context, id uint) (*entities.Benificiary, error) {
	return i.BeneficiaryRepo.GetBeneficiaryByID(ctx, id)
}

// GetCustomerBenificiaries implements ICustomerService.
func (i *IndividualCustomerService) GetCustomerBeneficiaries(ctx context.Context, id uint) ([]*entities.Benificiary, error) {
	return i.BeneficiaryRepo.GetBeneficiariesByCustomerID(ctx, id)
}

func (i *IndividualCustomerService) CreateBeneficiary(ctx context.Context, in CreateBenificiaryIn) (*entities.Benificiary, error) {
	// Validate input
	b := &entities.Benificiary{
		OwnerCustomerID: in.CustomerID,
		FirstName:       in.FirstName,
		LastName:        in.LastName,
		PhoneNumber:     in.PhoneNumber,
		Email:           in.Email,
		CreatedAt:       time.Now().UTC(),
	}

	if err := b.Validate(); err != nil {
		return nil, err
	}

	cl, err := i.Core.GetClientByExternalId(ctx, in.PhoneNumber)
	if err != nil {
		return nil, err
	}

	b.BeneficiaryCustomerID = cl.Id
	if err := i.BeneficiaryRepo.SaveBeneficiary(ctx, b); err != nil {
		return nil, err
	}

	return b, nil
}
