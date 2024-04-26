package repository

import (
	"context"

	"github.com/bekha-io/vaultonomy/domain/entities"
	"github.com/bekha-io/vaultonomy/domain/types"
)

type IIndividualCustomerRepository interface {
	GetByID(ctx context.Context, id types.CustomerID) (*entities.IndividualCustomer, error)
	GetBy(ctx context.Context, key string, value interface{}) (*entities.IndividualCustomer, error)
	Save(ctx context.Context, customer *entities.IndividualCustomer) error
}
