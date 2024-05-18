package accounts

import (
	"github.com/bekha-io/openbank/domain/types"
	"github.com/gin-gonic/gin"
)


func (ctrl *AccountsController) GetAccount(c *gin.Context) {
	accountId := c.Param("id")
    if len(accountId) < 12 {
        c.JSON(400, gin.H{"error": "len"})
        return
    }

	account, err := ctrl.AccountsService.GetAccountByID(c, types.AccountID(accountId))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
        return
	}

	c.JSON(200, account)
}