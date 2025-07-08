package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/repository"
)

// MemoryTransactionRepository implements the ITransactionRepository interface with in-memory storage.
type MemoryTransactionRepository struct {
	mu           sync.RWMutex
	transactions map[uint]*entities.Transaction
}

// NewMemoryTransactionRepository creates a new instance of MemoryTransactionRepository.
func NewMemoryTransactionRepository() *MemoryTransactionRepository {
	return &MemoryTransactionRepository{
		transactions: make(map[uint]*entities.Transaction),
	}
}

// GetByID retrieves a transaction by its ID.
func (m *MemoryTransactionRepository) GetByID(ctx context.Context, id uint) (*entities.Transaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tr, exists := m.transactions[id]
	if !exists {
		return nil, errors.New("transaction not found")
	}
	return tr, nil
}

// GetTransactionsByAccountID retrieves transactions by the account ID.
func (m *MemoryTransactionRepository) GetTransactions(ctx context.Context, in repository.GetTransactionsIn) ([]*entities.Transaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*entities.Transaction
	for _, tr := range m.transactions {
		// Проверяем принадлежность к аккаунту
		if tr.FromAccountId == in.AccountID || tr.ToAccountId == in.AccountID {
			// Проверяем фильтры по дате
			if !in.DateFrom.IsZero() && tr.CreatedAt.Before(in.DateFrom) {
				continue
			}
			if !in.DateTo.IsZero() && tr.CreatedAt.After(in.DateTo) {
				continue
			}
			result = append(result, tr)
		}
	}
	return result, nil
}

// SaveTransaction saves a transaction to the in-memory store.
func (m *MemoryTransactionRepository) SaveTransaction(ctx context.Context, tr *entities.Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.transactions[tr.ID] = tr
	return nil
}

// Ensure MemoryTransactionRepository implements the ITransactionRepository interface.
var _ repository.ITransactionRepository = (*MemoryTransactionRepository)(nil)
