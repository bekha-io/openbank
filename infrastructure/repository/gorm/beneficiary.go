package gormRepo

import (
	"context"
	"time"

	"github.com/bekha-io/openbank/domain/entities"
	"gorm.io/gorm"
)

type gormBeneficiary struct {
	ID                    uint      `gorm:"primaryKey"`
	OwnerCustomerID       uint
	BeneficiaryCustomerID uint
	FirstName             string
	LastName              string
	PhoneNumber           string
	Email                 string
	CreatedAt             time.Time
}

func (g *gormBeneficiary) toEntity() *entities.Benificiary {
	return &entities.Benificiary{
		ID:                    g.ID,
		OwnerCustomerID:       g.OwnerCustomerID,
		BeneficiaryCustomerID: g.BeneficiaryCustomerID,
		FirstName:             g.FirstName,
		LastName:              g.LastName,
		PhoneNumber:           g.PhoneNumber,
		Email:                 g.Email,
		CreatedAt:             g.CreatedAt,
	}
}

func fromEntity(e *entities.Benificiary) *gormBeneficiary {
	return &gormBeneficiary{
		ID:                    e.ID,
		OwnerCustomerID:       e.OwnerCustomerID,
		BeneficiaryCustomerID: e.BeneficiaryCustomerID,
		FirstName:             e.FirstName,
		LastName:              e.LastName,
		PhoneNumber:           e.PhoneNumber,
		Email:                 e.Email,
		CreatedAt:             e.CreatedAt,
	}
}


type PostgresBenificiaryRepository struct {
	db *gorm.DB
}

func NewPostgresBenificiaryRepository(db *gorm.DB) *PostgresBenificiaryRepository {
	return &PostgresBenificiaryRepository{db: db}
}

func (r *PostgresBenificiaryRepository) GetBeneficiaryByID(ctx context.Context, id uint) (*entities.Benificiary, error) {
	var b gormBeneficiary
	if err := r.db.WithContext(ctx).First(&b, id).Error; err != nil {
		return nil, err
	}
	return b.toEntity(), nil
}

func (r *PostgresBenificiaryRepository) GetBeneficiariesByCustomerID(ctx context.Context, customerId uint) ([]*entities.Benificiary, error) {
	var list []gormBeneficiary
	if err := r.db.WithContext(ctx).
		Where("owner_customer_id = ?", customerId).
		Find(&list).Error; err != nil {
		return nil, err
	}

	var result []*entities.Benificiary
	for _, g := range list {
		result = append(result, g.toEntity())
	}
	return result, nil
}

func (r *PostgresBenificiaryRepository) SaveBeneficiary(ctx context.Context, e *entities.Benificiary) error {
	g := fromEntity(e)
	return r.db.WithContext(ctx).Save(g).Error
}

func (r *PostgresBenificiaryRepository) DeleteBeneficiary(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&gormBeneficiary{}, id).Error
}
