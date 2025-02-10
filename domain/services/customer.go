package services

import (
	"context"

	"github.com/bekha-io/openbank/domain/dto"
	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/types"
)

type ICustomerService interface {
	CreateCustomer(ctx context.Context, in dto.CreateIndividualCustomerCommand) error
	GetCustomer(ctx context.Context, id types.CustomerID) (*entities.Customer, error)
	GetCustomerBy(ctx context.Context, key string, value interface{}) (*entities.Customer, error)
	SearchCustomersByPhoneNumber(ctx context.Context, phoneNumber string) ([]*entities.Customer, error)
	GetCustomerAccounts(ctx context.Context, customer entities.Customer) ([]*entities.Account, error)
}
