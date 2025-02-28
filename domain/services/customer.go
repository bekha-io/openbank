package services

import (
	"context"

	"github.com/bekha-io/openbank/domain/dto"
	"github.com/bekha-io/openbank/domain/entities"
)

type CreateBenificiaryIn struct {
	CustomerID  uint   `json:"customer_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type ICustomerService interface {
	CreateCustomer(ctx context.Context, in dto.CreateIndividualCustomerCommand) error
	GetCustomer(ctx context.Context, id uint) (*entities.Customer, error)
	GetCustomerByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.Customer, error)
	GetCustomerAccounts(ctx context.Context, customer entities.Customer) ([]*entities.Account, error)

	// Beneficiaries
	GetCustomerBeneficiaries(ctx context.Context, id uint) ([]*entities.Benificiary, error)
	GetBeneficiaryByID(ctx context.Context, id uint) (*entities.Benificiary, error)
	CreateBeneficiary(ctx context.Context, in CreateBenificiaryIn) (*entities.Benificiary, error)
}
