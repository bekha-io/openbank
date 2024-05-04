package entities

import "github.com/bekha-io/openbank/domain/types"

type Customer interface {
	Id() types.CustomerID
}
