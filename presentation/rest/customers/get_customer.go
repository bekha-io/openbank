package customers

import (
	"github.com/bekha-io/openbank/domain/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


func (ctrl *CustomerController) GetCustomer(c *gin.Context) {
	id := c.Param("id")

	uui, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	customer, err := ctrl.IndividualCustomerService.GetCustomer(c, types.CustomerID(uui))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, customer)
}