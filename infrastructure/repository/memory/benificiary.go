package memory

import (
	"context"
	"sync"

	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/repository"
	"github.com/bekha-io/openbank/domain/types/errs"
)

var _ repository.IBenificiaryRepository = (*MemoryBenificiaryRepository)(nil)

type MemoryBenificiaryRepository struct {
	m           sync.Mutex
	counter     uint
	data        map[uint]*entities.Benificiary
	customerMap map[uint][]*entities.Benificiary
}

func NewMemoryBeneficiaryRepository() *MemoryBenificiaryRepository {
	return &MemoryBenificiaryRepository{
		counter:     0,
		data:        make(map[uint]*entities.Benificiary),
		customerMap: make(map[uint][]*entities.Benificiary),
	}
}

// DeleteBenificiary implements repository.IBenificiaryRepository.
func (r *MemoryBenificiaryRepository) DeleteBeneficiary(ctx context.Context, id uint) error {
	r.m.Lock()
	defer r.m.Unlock()

	ben, ok := r.data[id]
	if !ok {
		// No-op
		return nil
	}

	for i, v := range r.customerMap[ben.OwnerCustomerID] {
		if v.ID == id {
			r.customerMap[ben.OwnerCustomerID] = append(r.customerMap[ben.OwnerCustomerID][:i], r.customerMap[ben.OwnerCustomerID][i+1:]...)
			break
		}
	}

	delete(r.data, id)
	return nil
}

// GetBenificiariesByCustomerID implements repository.IBenificiaryRepository.
func (m *MemoryBenificiaryRepository) GetBeneficiariesByCustomerID(ctx context.Context, customerId uint) ([]*entities.Benificiary, error) {
	m.m.Lock()
	defer m.m.Unlock()

	if bens, ok := m.customerMap[customerId]; ok {
		return bens, nil
	}

	return []*entities.Benificiary{}, nil
}

// GetBenificiaryByID implements repository.IBenificiaryRepository.
func (m *MemoryBenificiaryRepository) GetBeneficiaryByID(ctx context.Context, id uint) (*entities.Benificiary, error) {
	m.m.Lock()
	defer m.m.Unlock()

	ben, ok := m.data[id]
	if ok {
		return ben, nil
	}

	return nil, errs.ErrBenificiaryNotFound
}

// SaveBenificiary implements repository.IBenificiaryRepository.
func (m *MemoryBenificiaryRepository) SaveBeneficiary(ctx context.Context, benificiary *entities.Benificiary) error {
	m.m.Lock()
	defer m.m.Unlock()

	benificiary.ID = m.counter
	m.counter++
	m.data[benificiary.ID] = benificiary

	if bens, ok := m.customerMap[benificiary.OwnerCustomerID]; ok {
		m.customerMap[benificiary.OwnerCustomerID] = append(bens, benificiary)
	} else {
		m.customerMap[benificiary.OwnerCustomerID] = []*entities.Benificiary{benificiary}
	}

	return nil
}
