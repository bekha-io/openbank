package accounts

import (
	"net/http"

	"github.com/bekha-io/openbank/domain/types"
	"github.com/gin-gonic/gin"
)

func (ctrl *AccountsController) SearchAccounts(c *gin.Context) {
	id := c.Query("id")

	if len(id) < 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}

	accounts, err := ctrl.AccountsService.GetAccountsLike(c, types.AccountID(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, accounts)
}
