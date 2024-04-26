package customers

import (
	"net/http"

	"github.com/bekha-io/vaultonomy/domain/dto"
	"github.com/gin-gonic/gin"
)

func (ctrl *CustomerController) CreateCustomer(c *gin.Context) {
	type req struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
	}

	var in req

	err := c.BindJSON(&in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ctrl.IndividualCustomerService.CreateCustomer(c.Request.Context(),
		dto.CreateIndividualCustomerCommand{
			PhoneNumber: in.PhoneNumber,
		})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, "user created")
}
