package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/bekha-io/openbank/domain/services"
	"github.com/bekha-io/openbank/domain/types/errs"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthController struct {
	AuthorizationService services.IAuthorizationService
	EmployeeService      services.IEmployeeService
}

var (
	// Only for testing purposes. Should not be used in production
	jwtSigningKey = []byte("openbankjwtsigningkey")
)

func NewAuthController(authSvc services.IAuthorizationService, emplSvc services.IEmployeeService) *AuthController {
	return &AuthController{
		EmployeeService:      emplSvc,
		AuthorizationService: authSvc,
	}
}

func (ctrl *AuthController) generateJwtToken(expiresAt time.Time, values map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{
		"iss": "openbank-api-server",
		"iat": time.Now().UTC().Unix(),
		"exp": expiresAt.Unix(),
	}
	for k, v := range values {
		claims[k] = v
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString(jwtSigningKey)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (ctrl *AuthController) parseToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return jwtSigningKey, nil
	}, jwt.WithExpirationRequired(), jwt.WithIssuedAt())
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("token is not valid")
	}

	return parsedToken.Claims.(jwt.MapClaims), nil
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

	token, err := ctrl.generateJwtToken(time.Now().UTC().Add(time.Hour*12), map[string]interface{}{
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

		claims, err := ctrl.parseToken(bearerToken)
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
