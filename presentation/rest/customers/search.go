package customers

import "github.com/gin-gonic/gin"


func (ctrl *CustomerController) SearchCustomers(c *gin.Context) {
	phoneNumber := c.Query("phone_number")
	if len(phoneNumber) < 4 {
        c.JSON(400, gin.H{"error": "missing phone number"})
        return
    }

	customers, err := ctrl.IndividualCustomerService.SearchCustomersByPhoneNumber(c, phoneNumber)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
        return
	}
	c.JSON(200, customers)
}