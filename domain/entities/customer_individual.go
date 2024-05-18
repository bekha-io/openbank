package entities

import (
	"fmt"
	"time"

	"github.com/bekha-io/openbank/domain/types"
	"github.com/google/uuid"
)

type IndividualCustomer struct {
	ID          types.CustomerID `json:"id"`
	PhoneNumber string `json:"phone_number"` // without a leading +

	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`

	Passport *Passport `json:"passport"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
