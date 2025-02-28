package repository

import (
	"context"

	"github.com/bekha-io/openbank/domain/entities"
)

type IBenificiaryRepository interface {
	GetBeneficiaryByID(ctx context.Context, id uint) (*entities.Benificiary, error)
	GetBeneficiariesByCustomerID(ctx context.Context, customerId uint) ([]*entities.Benificiary, error)
	SaveBeneficiary(ctx context.Context, benificiary *entities.Benificiary) error
	DeleteBeneficiary(ctx context.Context, id uint) error
}
