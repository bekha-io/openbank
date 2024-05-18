package types

import (
	"encoding/json"

	"github.com/google/uuid"
)

type EmployeeID uuid.UUID

func (e EmployeeID) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e EmployeeID) String() string {
	return uuid.UUID(e).String()
}

func (e EmployeeID) UUID() uuid.UUID {
	return uuid.UUID(e)
}