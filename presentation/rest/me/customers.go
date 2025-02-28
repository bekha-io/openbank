package me

import "github.com/gin-gonic/gin"

func (ctrl *Controller) GetCustomerByPhoneNumber(c *gin.Context) {
	phoneNumber := c.Param("phoneNumber")
	if len(phoneNumber) < 10 {
		c.JSON(400, gin.H{"error": "missing or invalid phone number"})
		return
	}

	customer, err := ctrl.CustomersService.GetCustomerByPhoneNumber(c, phoneNumber)
	if err != nil {
		handleError(c, 404, err)
		return
	}

	c.JSON(200, customer)
}
