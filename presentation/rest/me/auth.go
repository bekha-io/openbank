package me

import (
	"errors"
	"strings"
	"time"

	"github.com/bekha-io/openbank/domain/types/errs"
	"github.com/bekha-io/openbank/presentation/rest/utils"
	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) CustomerSignIn(c *gin.Context) {
	type req struct {
		// No authentication required. ONLY FOR FAST INTERFACE TESTING PURPOSES
		PhoneNumber string `json:"phone_number" binding:"required"`
	}

	var in req
	if err := c.BindJSON(&in); err != nil {
		handleError(c, 400, err)
		return
	}
	
	customer, err := ctrl.CustomersService.GetCustomerBy(c, "phone_number", in.PhoneNumber)
	if err != nil {
		handleError(c, 404, err)
		return
	}

	token, err := utils.GenerateJwtToken(time.Now().UTC().Add(time.Hour * 12), map[string]interface{}{
		"customer_id": customer.ID,
	})
	if err != nil {
		handleError(c, 404, err)
	}

	c.JSON(200, gin.H{"customer": customer, "token": token})
}


func (ctrl *Controller) CustomerAuthenticateMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
        bearerToken := strings.Replace(authHeader, "Bearer ", "", 1)

        claims, err := utils.ParseToken(bearerToken)
        if err!= nil {
            ctx.JSON(401, gin.H{"error": errors.Join(errs.ErrNotAuthenticated, err).Error()})
            ctx.Abort()
            return
        }

        customerId, okId := claims["customer_id"]
        if!okId {
            ctx.JSON(401, gin.H{"error": errs.ErrNotAuthenticated.Error()})
            ctx.Abort()
            return
        }

        ctx.Set("customerId", customerId)
	}
}