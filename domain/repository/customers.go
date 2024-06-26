package repository

import (
	"context"

	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/types"
)

type IIndividualCustomerRepository interface {
	GetByID(ctx context.Context, id types.CustomerID) (*entities.IndividualCustomer, error)
	GetBy(ctx context.Context, key string, value interface{}) (*entities.IndividualCustomer, error)
	GetManyIDLike(ctx context.Context, id types.Currency) ([]*entities.IndividualCustomer, error)
	GetManyPhoneNumberLike(ctx context.Context, phoneNumber string) ([]*entities.IndividualCustomer, error)
	Save(ctx context.Context, customer *entities.IndividualCustomer) error
}
