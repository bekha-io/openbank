package permissions

var (
	CanCrudRole   = NewPermission("roles", CRUD)
	CanCreateRole = NewPermission("roles", Create)
	CanReadRole   = NewPermission("roles", Read)
	CanUpdateRole = NewPermission("roles", Update)
	CanDeleteRole = NewPermission("roles", Delete)
)

var RolePermissionSet = []Permission{
	CanCrudRole, CanCreateRole, CanReadRole, CanUpdateRole, CanDeleteRole,
}
