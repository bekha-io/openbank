package repository

import (
	"context"

	"github.com/bekha-io/openbank/domain/entities"
)

type IAccountRepository interface {
	GetByID(ctx context.Context, id uint) (*entities.Account, error)
	GetBy(ctx context.Context, key string, value interface{}) (*entities.Account, error)
	GetManyIdLike(ctx context.Context, id uint) ([]*entities.Account, error)
	GetManyBy(ctx context.Context, filters ...Filter) ([]*entities.Account, error)
	Save(ctx context.Context, account *entities.Account) error
}
