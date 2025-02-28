package me

import (
	"strconv"

	"github.com/bekha-io/openbank/domain/services"
	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) GetBeneficiaryByID(c *gin.Context) {
	benificiaryId := c.Param("id")

	beneficiaryInt, err := strconv.Atoi(benificiaryId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	b, err := ctrl.CustomersService.GetBeneficiaryByID(c, uint(beneficiaryInt))
	if err != nil {
		handleError(c, 400, err)
		return
	}

	c.JSON(200, b)
}

func (ctrl *Controller) GetCustomerBeneficiaries(c *gin.Context) {
	customerId := c.GetFloat64("customerId")
	if customerId == 0 {
		c.JSON(400, gin.H{"error": "missing customer ID"})
		return
	}

	bs, err := ctrl.CustomersService.GetCustomerBeneficiaries(c, uint(customerId))
	if err != nil {
		handleError(c, 400, err)
		return
	}

	c.JSON(200, bs)
}

func (ctrl *Controller) CreateBeneficiary(c *gin.Context) {
	customerId := c.GetFloat64("customerId")
	if customerId == 0 {
		c.JSON(400, gin.H{"error": "missing customer ID"})
		return
	}

	var in = struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		PhoneNumber string `json:"phone_number"`
		Email       string `json:"email"`
	}{}

	if err := c.BindJSON(&in); err != nil {
		handleError(c, 400, err)
		return
	}

	b, err := ctrl.CustomersService.CreateBeneficiary(c, services.CreateBenificiaryIn{
		CustomerID:  uint(customerId),
		FirstName:   in.FirstName,
		LastName:    in.LastName,
		PhoneNumber: in.PhoneNumber,
		Email:       in.Email,
	})
	if err != nil {
		handleError(c, 400, err)
		return
	}

	c.JSON(200, b)
}
