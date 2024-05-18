package types

import (
	"encoding/json"

	"github.com/google/uuid"
)

type BranchID uuid.UUID

func (b BranchID) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

func (b BranchID) String() string {
	return uuid.UUID(b).String()
}

func (b BranchID) UUID() uuid.UUID {
	return uuid.UUID(b)
}