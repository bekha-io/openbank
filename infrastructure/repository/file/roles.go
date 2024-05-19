package file

import (
	"io"
	"os"
	"path"
	"strings"

	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/entities/permissions"
	"github.com/bekha-io/openbank/domain/repository"
)

var _ repository.IRolesRepository = (*FileRolesRepository)(nil)

type FileRolesRepository struct {
	Dir string
}

// GetRoleByName implements repository.IRolesRepository.
func (f *FileRolesRepository) GetRoleByName(name string) (*entities.Role, error) {
	file, err := os.Open(path.Join(f.Dir, name))
	if err != nil {
		return nil, err
	}

	var role = &entities.Role{}
	var perms = []permissions.Permission{}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	for i, value := range lines {
		if i == 0 {
			role.Name = value
		}

		if i == 1 {
			role.Description = value
		}

		if len(value) > 0 {
			permData := strings.Split(value, ":")
			perms = append(perms, permissions.Permission{
				Object:      permData[0],
				Action:      permData[1],
				Description: permData[2],
			})
		}
	}
	role.Permissions = perms

	return role, nil
}

// SaveRole implements repository.IRolesRepository.
func (f *FileRolesRepository) SaveRole(role *entities.Role) error {
	file, err := os.Create(path.Join(f.Dir, role.Name))
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString(role.Name + "\n")
	file.WriteString(role.Description + "\n")

	for _, perm := range role.Permissions {
		file.WriteString(perm.Object + ":" + perm.Action + ":" + perm.Description + "\n")
	}

	file.Sync()
	return nil
}
