package accounts

import (
	"net/http"

	"github.com/bekha-io/vaultonomy/domain/dto"
	"github.com/bekha-io/vaultonomy/domain/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (ctrl *AccountsController) CreateAccount(c *gin.Context) {
	type req struct {
		CustomerID uuid.UUID `json:"customer_id"`
		Currency   string    `json:"currency"`
	}

	var in req
	err := c.BindJSON(&in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ctrl.AccountsService.CreateAccount(c, dto.CreateAccountCommand{
		CustomerID: types.CustomerID(in.CustomerID), Currency: types.Currency(in.Currency)})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, "account created")
}
