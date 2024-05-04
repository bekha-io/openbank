package customers

import (
	"github.com/bekha-io/openbank/domain/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (ctrl *CustomerController) GetCustomerAccounts(c *gin.Context) {
	customerId := c.Param("id")

	customerUuid, err := uuid.Parse(customerId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	customer, err := ctrl.IndividualCustomerService.GetCustomer(c, types.CustomerID(customerUuid))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	accounts, err := ctrl.IndividualCustomerService.GetCustomerAccounts(c, customer)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, accounts)
}
