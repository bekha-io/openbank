package entities

import "github.com/bekha-io/vaultonomy/domain/types"

type Customer interface {
	Id() types.CustomerID
}
