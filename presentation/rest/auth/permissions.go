package auth

import (
	"errors"

	"github.com/bekha-io/openbank/domain/entities/permissions"
	"github.com/gin-gonic/gin"
)

func ForbiddenJSON(c *gin.Context, err error) {
	c.JSON(403, gin.H{"error": err})
	c.Abort()
	return
}

func (ctrl *AuthController) IfHasPermissions(fc gin.HandlerFunc, perms ...permissions.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleNames := c.GetStringSlice("roles")

		var permsMap = map[permissions.Permission]bool{}
		for _, perm := range perms {
			permsMap[perm] = false
		}

		for _, roleName := range roleNames {
			role, err := ctrl.AuthorizationService.GetRoleByName(c, roleName)
			if err != nil {
				ForbiddenJSON(c, err)
			}

			// Role permissions
			for _, rolePerm := range role.Permissions {
				for _, perm := range perms {
					if rolePerm.Is(perm) {
						permsMap[perm] = true
					}
				}
			}
		}

		for k, v := range permsMap {
			if !v {
				ForbiddenJSON(c, errors.New(k.String()))
				return
			}
		}

		c.Next()
	}
}
