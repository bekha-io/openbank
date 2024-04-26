package dto

import "errors"

type CreateIndividualCustomerCommand struct {
	PhoneNumber string
}

func (c CreateIndividualCustomerCommand) Validate() error {
	if c.PhoneNumber == "" {
		return errors.New("phone number must not be empty")
	}
	return nil
}