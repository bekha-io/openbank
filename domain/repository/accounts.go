package repository

import (
	"context"

	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/types"
)

type IAccountRepository interface {
	GetByID(ctx context.Context, id types.AccountID) (*entities.Account, error)
	GetBy(ctx context.Context, key string, value interface{}) (*entities.Account, error)
	GetManyIdLike(ctx context.Context, id types.AccountID) ([]*entities.Account, error)
	GetManyBy(ctx context.Context, filters ...Filter) ([]*entities.Account, error)
	Save(ctx context.Context, account *entities.Account) error
}