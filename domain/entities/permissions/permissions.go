package permissions

const (
	Create = "create"
	Update = "update"
	Delete = "delete"
	Read   = "read"
	CRUD   = "crud"
)

type Permission struct {
	Object      string
	Action      string
	Description string
}

func (p Permission) Is(pm Permission) bool {
	if p.Object == pm.Object && p.Action == CRUD {
		return true
	}
	return p.Object == pm.Object && p.Action == pm.Action
}

func (p Permission) String() string {
	return p.Object + ":" + p.Action
}

func NewPermission(object string, action string) Permission {
	return Permission{
		Object: object,
		Action: action,
	}
}

func CreatePermissionSet(object string) []Permission {
	return []Permission{
		NewPermission(object, Create),
		NewPermission(object, Update),
		NewPermission(object, Delete),
		NewPermission(object, Read),
		NewPermission(object, CRUD),
	}
}
