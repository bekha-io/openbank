package repository

import (
	"context"
	"time"

	"github.com/bekha-io/openbank/domain/entities"
)

type GetTransactionsIn struct {
	AccountID uint
	DateFrom time.Time
	DateTo   time.Time
}

type ITransactionRepository interface {
	GetByID(ctx context.Context, id uint) (*entities.Transaction, error)
	GetTransactions(ctx context.Context, in GetTransactionsIn) ([]*entities.Transaction, error)
	SaveTransaction(ctx context.Context, tr *entities.Transaction) error
}
