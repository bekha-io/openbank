package repository

import "github.com/bekha-io/openbank/domain/entities"

type IRolesRepository interface {
	GetRoleByName(name string) (*entities.Role, error)
	SaveRole(role *entities.Role) error
}
