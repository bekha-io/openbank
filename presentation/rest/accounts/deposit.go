package accounts

import (
	"net/http"

	"github.com/bekha-io/vaultonomy/domain/dto"
	"github.com/bekha-io/vaultonomy/domain/types"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func (ctrl *AccountsController) Deposit(c *gin.Context) {
	accountId := c.Param("id")
	type req struct {
		Amount   decimal.Decimal `json:"amount" binding:"required"`
		Currency types.Currency  `json:"currency" binding:"required"`
		Comment  string          `json:"comment"`
	}

	var in req
	err := c.BindJSON(&in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := ctrl.AccountsService.GetAccountByID(c, types.AccountID(accountId))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := ctrl.AccountsService.Deposit(c, dto.DepositCommand{
		Account: account,
		Comment: in.Comment,
		Money:   *types.NewMoney(in.Amount, in.Currency),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}
