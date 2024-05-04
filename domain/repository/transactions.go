package repository

import (
	"context"

	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/types"
)

type ITransactionRepository interface {
	GetByID(ctx context.Context, id types.TransactionID) (*entities.Transaction, error)
	GetBy(ctx context.Context, key string, value interface{}) (*entities.Transaction, error)
	GetManyBy(ctx context.Context, filters ...Filter) ([]*entities.Transaction, error)
	Save(ctx context.Context, transaction *entities.Transaction) error
}
