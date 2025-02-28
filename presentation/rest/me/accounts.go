package me

import (
	"net/http"
	"strconv"

	"github.com/bekha-io/openbank/domain/services"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func (ctrl *Controller) GetAccounts(c *gin.Context) {
	customerId := c.GetFloat64("customerId")

	customer, err := ctrl.CustomersService.GetCustomer(c, uint(customerId))
	if err != nil {
		handleError(c, 400, err)
		return
	}

	accounts, err := ctrl.CustomersService.GetCustomerAccounts(c, *customer)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"accounts": accounts})
}

func (ctrl *Controller) GetAccountTransactions(c *gin.Context) {
	accountIdStr := c.Param("id")

	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		handleError(c, 400, err)
		return
	}

	customerId := c.GetFloat64("customerId")
	if customerId == 0 {
		c.JSON(403, gin.H{"error": "missing customer ID"})
		return
	}

	account, err := ctrl.AccountsService.GetAccountByID(c, uint(accountId))
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error(), "account_id": accountId})
		return
	}

	if account.CustomerID != uint(customerId) {
		c.JSON(403, gin.H{"error": "customer is not the account's owner"})
		return
	}

	transactions, err := ctrl.AccountsService.GetAccountTransactions(c, uint(accountId))
	if err != nil {
		handleError(c, 400, err)
		return
	}

	c.JSON(200, gin.H{"transactions": transactions})
}

func (ctrl *Controller) TransferMoney(c *gin.Context) {
	var in = struct {
		FromAccountId   uint            `json:"from_account_id"`
		ToBeneficiaryId uint            `json:"to_beneficiary_id"`
		Amount          decimal.Decimal `json:"amount"`
		Comment         string          `json:"comment"`
	}{}

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Perform the transfer
	transaction, err := ctrl.AccountsService.Transfer(c, services.TransferIn{
		FromAccountId:   in.FromAccountId,
		ToBeneficiaryId: in.ToBeneficiaryId,
		Amount:          in.Amount,
		Comment:         in.Comment,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "transaction": transaction})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}
