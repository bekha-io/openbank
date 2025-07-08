package gormRepo

import (
	"context"
	"time"

	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/repository"
	"github.com/bekha-io/openbank/domain/types"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type gormTransaction struct {
	ID            uint `gorm:"primaryKey"`
	ExternalID    string
	FromAccountId uint
	ToAccountId   uint
	Category      string
	Status        string
	StatusReason  string
	Comment       string
	Amount        decimal.Decimal `gorm:"type:numeric"`
	CreatedAt     time.Time
}

func (g *gormTransaction) toEntity() *entities.Transaction {
	return &entities.Transaction{
		ID:            g.ID,
		ExternalID:    g.ExternalID,
		FromAccountId: g.FromAccountId,
		ToAccountId:   g.ToAccountId,
		Category:      types.TransactionCategory(g.Category),
		Status:        entities.TransactionStatus(g.Status),
		StatusReason:  g.StatusReason,
		Comment:       g.Comment,
		Amount:        g.Amount,
		CreatedAt:     g.CreatedAt,
	}
}

func (g *gormTransaction) fromEntity(e *entities.Transaction) {
	g.ID = e.ID
	g.ExternalID = e.ExternalID
	g.FromAccountId = e.FromAccountId
	g.ToAccountId = e.ToAccountId
	g.Category = string(e.Category)
	g.Status = string(e.Status)
	g.StatusReason = e.StatusReason
	g.Comment = e.Comment
	g.Amount = e.Amount
	g.CreatedAt = e.CreatedAt
}

type GormTransactionRepository struct {
	db *gorm.DB
}

func NewPostgresTransactionRepository(db *gorm.DB) *GormTransactionRepository {
	return &GormTransactionRepository{db: db}
}

var _ repository.ITransactionRepository = (*GormTransactionRepository)(nil)

func (r *GormTransactionRepository) GetByID(ctx context.Context, id uint) (*entities.Transaction, error) {
	var g gormTransaction
	if err := r.db.WithContext(ctx).First(&g, id).Error; err != nil {
		return nil, err
	}
	return g.toEntity(), nil
}

func (r *GormTransactionRepository) GetTransactions(ctx context.Context, in repository.GetTransactionsIn) ([]*entities.Transaction, error) {
	query := r.db.WithContext(ctx).Where("from_account_id = ? OR to_account_id = ?", in.AccountID, in.AccountID)

	// Добавляем фильтры по дате если они указаны
	if !in.DateFrom.IsZero() {
		query = query.Where("created_at >= ?", in.DateFrom)
	}
	if !in.DateTo.IsZero() {
		query = query.Where("created_at <= ?", in.DateTo)
	}

	var list []gormTransaction
	if err := query.Find(&list).Error; err != nil {
		return nil, err
	}

	var result []*entities.Transaction
	for _, g := range list {
		result = append(result, g.toEntity())
	}
	return result, nil
}

func (r *GormTransactionRepository) SaveTransaction(ctx context.Context, e *entities.Transaction) error {
	g := gormTransaction{}
	g.fromEntity(e)
	return r.db.WithContext(ctx).Save(&g).Error
}
