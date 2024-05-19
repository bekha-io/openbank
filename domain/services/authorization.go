package services

import (
	"context"

	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/entities/permissions"
	"github.com/bekha-io/openbank/domain/repository"
)

type IAuthorizationService interface {
	GetAllPermissions(ctx context.Context) []permissions.Permission
	GetRoleByName(ctx context.Context, name string) (*entities.Role, error)
}

var _ IAuthorizationService = (*AuthorizationService)(nil)

type AuthorizationService struct {
	RolesRepository repository.IRolesRepository
}

func NewAuthorizationService(rolesRepository repository.IRolesRepository) *AuthorizationService {
	return &AuthorizationService{
		RolesRepository: rolesRepository,
	}
}

func (s *AuthorizationService) GetAllPermissions(ctx context.Context) []permissions.Permission {
	var perms []permissions.Permission
	perms = append(perms, permissions.RolePermissionSet...)
	perms = append(perms, permissions.EmployeePermissionSet...)
	perms = append(perms, permissions.LoanProductPermissionSet...)
	return perms
}

func (s *AuthorizationService) GetRoleByName(ctx context.Context, name string) (*entities.Role, error) {
	return s.RolesRepository.GetRoleByName(name)
}
