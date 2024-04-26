package types

import (
	"encoding/json"

	"github.com/biter777/countries"
	"github.com/google/uuid"
)

type CustomerID uuid.UUID

func (c CustomerID) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c CustomerID) String() string {
	return uuid.UUID(c).String()
}

func (c CustomerID) UUID() uuid.UUID {
	return uuid.UUID(c)
}

type Country countries.Country
