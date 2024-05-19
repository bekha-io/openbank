package repository

import (
	"context"

	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/types"
)

type ILoanRepository interface {
	SaveLoanProduct(ctx context.Context, lp *entities.LoanProduct) error
	GetLoanProductByID(ctx context.Context, id types.LoanProductID) (*entities.LoanProduct, error)
	GetAllLoanProducts(ctx context.Context) ([]*entities.LoanProduct, error)
}
