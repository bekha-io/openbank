package entities

import (
	"time"

	"github.com/bekha-io/vaultonomy/domain/types"
)

type Passport struct {
	Country types.Country

	Serial string
	Number string

	IssuedAt  time.Time
	ExpiresAt time.Time

	Nationality types.Country

	CreatedAt time.Time
	UpdatedAt time.Time
}
