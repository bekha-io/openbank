package entities

import (
	"fmt"
	"time"
)

type Benificiary struct {
	ID                    uint      `json:"id"`
	OwnerCustomerID       uint      `json:"owner_customer_id"`
	BeneficiaryCustomerID uint      `json:"beneficiary_customer_id"`
	FirstName             string    `json:"first_name"`
	LastName              string    `json:"last_name"`
	PhoneNumber           string    `json:"phone_number"`
	Email                 string    `json:"email"`
	CreatedAt             time.Time `json:"created_at"`
}

func (b *Benificiary) Validate() error {
	if b.FirstName == "" || b.LastName == "" {
		return fmt.Errorf("first_name and last_name cannot be empty")
	}

	if b.PhoneNumber == "" {
		return fmt.Errorf("invalid phone number format")
	}

	if b.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	return nil
}
