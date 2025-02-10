package repository

import (
	"context"

	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/types"
)

type IIndividualCustomerRepository interface {
	GetByID(ctx context.Context, id types.CustomerID) (*entities.Customer, error)
	GetBy(ctx context.Context, key string, value interface{}) (*entities.Customer, error)
	GetManyIDLike(ctx context.Context, id types.Currency) ([]*entities.Customer, error)
	GetManyPhoneNumberLike(ctx context.Context, phoneNumber string) ([]*entities.Customer, error)
	Save(ctx context.Context, customer *entities.Customer) error
}
