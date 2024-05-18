package repository

import (
	"context"

	"github.com/bekha-io/openbank/domain/entities"
)

type IEmployeeRepository interface {
	GetByEmail(ctx context.Context, email string) (*entities.Employee, error)
	GetManyLike(ctx context.Context, lookup string) ([]*entities.Employee, error)
	Save(ctx context.Context, employee *entities.Employee) error
}