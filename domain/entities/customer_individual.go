package entities

import (
	"fmt"
)

type Customer struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phone_number"` // without a leading +

	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`

	Passport *Passport `json:"passport"`

}

func (c Customer) FullName() string {
	return fmt.Sprintf("%v %v %v", c.LastName, c.FirstName, c.MiddleName)
}
