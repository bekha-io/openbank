package repository

import (
	"context"

	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/types"
)

// type Result[T any] struct {
// 	Rows  T
// 	Error error
// 	Tx    Tx
// }

// type NoRowResult = Result[any]  // NoRowResult means no rows (Result.Rows) should be expected

type IAccountRepository interface {
	GetByID(ctx context.Context, id types.AccountID) (*entities.Account, error)
	GetBy(ctx context.Context, key string, value interface{}) (*entities.Account, error)
	GetManyBy(ctx context.Context, filters ...Filter) ([]*entities.Account, error)
	Save(ctx context.Context, account *entities.Account) error
}
