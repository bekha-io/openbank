package entities

import "github.com/bekha-io/openbank/domain/entities/permissions"

type Role struct {
	Name        string
	Description string
	Permissions []permissions.Permission
}
