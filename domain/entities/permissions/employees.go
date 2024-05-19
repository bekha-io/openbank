package permissions

var (
	CanCrudEmployee   = NewPermission("employees", CRUD)
	CanCreateEmployee = NewPermission("employees", Create)
	CanUpdateEmployee = NewPermission("employees", Update)
	CanDeleteEmployee = NewPermission("employees", Delete)
	CanReadEmployee   = NewPermission("employees", Read)
)

var EmployeePermissionSet = []Permission{
	CanCrudEmployee, CanCreateEmployee, CanUpdateEmployee,
	CanReadEmployee, CanDeleteEmployee,
}
