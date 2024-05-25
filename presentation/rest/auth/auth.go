package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/bekha-io/openbank/domain/services"
	"github.com/bekha-io/openbank/domain/types/errs"
	"github.com/bekha-io/openbank/presentation/rest/utils"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthorizationService services.IAuthorizationService
	EmployeeService      services.IEmployeeService
}

func NewAuthController(authSvc services.IAuthorizationService, emplSvc services.IEmployeeService) *AuthController {
	return &AuthController{
		EmployeeService:      emplSvc,
		AuthorizationService: authSvc,
	}
}

func (ctrl *AuthController) EmployeeSignIn(c *gin.Context) {
	type req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var in req
	c.BindJSON(&in)

	employee, err := ctrl.EmployeeService.Authenticate(c, in.Email, in.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": errors.Join(errs.ErrNotAuthenticated, err).Error()})
		return
	}

	token, err := utils.GenerateJwtToken(time.Now().UTC().Add(time.Hour*12), map[string]interface{}{
		"employee_id": employee.ID,
		"roles":       employee.Roles,
	})
	if err != nil {
		c.JSON(401, gin.H{"error": errors.Join(errs.ErrNotAuthenticated, err).Error()})
		return
	}

	c.JSON(200, gin.H{"token": token, "employee": employee})
}

func (ctrl *AuthController) EmployeeAuthenticateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		bearerToken := strings.Replace(authHeader, "Bearer ", "", 1)

		claims, err := utils.ParseToken(bearerToken)
		if err != nil {
			c.JSON(401, gin.H{"error": errors.Join(errs.ErrNotAuthenticated, err).Error()})
			c.Abort()
			return
		}

		employeeId, okId := claims["employee_id"]
		roles, okRoles := claims["roles"]
		if !okId || !okRoles {
			c.JSON(401, gin.H{"error": errs.ErrNotAuthenticated.Error()})
			c.Abort()
			return
		}

		c.Set("employeeId", employeeId)
		c.Set("roles", roles)
	}
}
