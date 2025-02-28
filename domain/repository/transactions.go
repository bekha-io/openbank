package repository

import (
	"context"

	"github.com/bekha-io/openbank/domain/entities"
)

type ITransactionRepository interface {
	GetByID(ctx context.Context, id uint) (*entities.Transaction, error)
	GetTransactionsByAccountID(ctx context.Context, id uint) ([]*entities.Transaction, error)
	SaveTransaction(ctx context.Context, tr *entities.Transaction) error
}
