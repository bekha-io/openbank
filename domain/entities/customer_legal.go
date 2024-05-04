package entities

import (
	"time"

	"github.com/bekha-io/openbank/domain/types"
)

type LegalCustomer struct {
	ID        types.CustomerID
	LegalName string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (l LegalCustomer) Id() types.CustomerID {
	return l.ID
}

func (l LegalCustomer) FullName() string {
	return l.LegalName
}
