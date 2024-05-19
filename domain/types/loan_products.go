package types

import (
	"encoding/json"

	"github.com/google/uuid"
)

type LoanProductID uuid.UUID

func (e LoanProductID) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e LoanProductID) String() string {
	return uuid.UUID(e).String()
}

func (e LoanProductID) UUID() uuid.UUID {
	return uuid.UUID(e)
}

type LoanType string

const (
	LoanTypeAnnuity LoanType = "annuity"
)

func (e LoanType) IsValid() bool {
	switch e {
	case LoanTypeAnnuity:
		return true
	default:
		return false
	}
}
