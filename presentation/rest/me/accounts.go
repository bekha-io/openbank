package me

import (
	"github.com/bekha-io/openbank/domain/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (ctrl *Controller) GetAccounts(c *gin.Context) {
	customerId := c.GetString("customerId")

	customerUuid, err := uuid.Parse(customerId)
	if err != nil {
		handleError(c, 400, err)
		return
	}

	customer, err := ctrl.CustomersService.GetCustomer(c, types.CustomerID(customerUuid))
	if err != nil {
		handleError(c, 400, err)
		return
	}

	accounts, err := ctrl.CustomersService.GetCustomerAccounts(c, customer)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"accounts": accounts})
}
