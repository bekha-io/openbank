package entities

import (
	"fmt"
	"time"

	"github.com/bekha-io/openbank/domain/types"
	"github.com/google/uuid"
)

type IndividualCustomer struct {
	ID          types.CustomerID
	PhoneNumber string // without a leading +

	FirstName  string
	LastName   string
	MiddleName string

	Passport *Passport

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewIndividualCustomer(phoneNumber string) *IndividualCustomer {
	return &IndividualCustomer{
		ID:          types.CustomerID(uuid.New()),
		PhoneNumber: phoneNumber,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

func (c *IndividualCustomer) Id() types.CustomerID {
	return c.ID
}

func (c IndividualCustomer) FullName() string {
	return fmt.Sprintf("%v %v %v", c.LastName, c.FirstName, c.MiddleName)
}
