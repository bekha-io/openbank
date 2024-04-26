package services

import (
	"context"

	"github.com/bekha-io/vaultonomy/domain/dto"
	"github.com/bekha-io/vaultonomy/domain/entities"
	"github.com/bekha-io/vaultonomy/domain/types"
)

type IIndividualCustomerService interface {
	CreateCustomer(ctx context.Context, in dto.CreateIndividualCustomerCommand) error
	GetCustomer(ctx context.Context, id types.CustomerID) (*entities.IndividualCustomer, error)
	GetCustomerBy(ctx context.Context, key string, value interface{}) (*entities.IndividualCustomer, error)

	GetCustomerAccounts(ctx context.Context, customer entities.Customer) ([]*entities.Account, error)
}

type ILegalCustomerService interface {
	CreateCustomer(ctx context.Context, in dto.CreateLegalCustomerCommand) error
	GetCustomer(ctx context.Context, id types.CustomerID) (*entities.Customer, error)
}
