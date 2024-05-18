package entities

import "github.com/bekha-io/openbank/domain/types"

type Branch struct {
	ID   types.BranchID `json:"id"`
	Name string         `json:"name"`
}