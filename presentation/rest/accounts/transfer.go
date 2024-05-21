package accounts

import (
	"fmt"

	"github.com/bekha-io/openbank/domain/dto"
	"github.com/bekha-io/openbank/domain/types"
	"github.com/bekha-io/openbank/domain/types/errs"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func (ctrl *AccountsController) Transfer(c *gin.Context) {
	type req struct {
		OriginAccountID      string          `json:"origin_account_id"`
		DestinationAccountID string          `json:"destination_account_id"`
		Amount               decimal.Decimal `json:"amount"`
		Currency             types.Currency  `json:"currency"`
	}

	var in req
	if err := c.BindJSON(&in); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}	

	originAcc, err := ctrl.AccountsService.GetAccountByID(c, types.AccountID(in.OriginAccountID))
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error(), "account_id": in.OriginAccountID})
		return
	}

	destAcc, err := ctrl.AccountsService.GetAccountByID(c, types.AccountID(in.DestinationAccountID))
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error(), "account_id": in.DestinationAccountID})
		return
	}

	if originAcc.Balance.Currency != destAcc.Balance.Currency {
		c.JSON(400, gin.H{"error": errs.ErrAccountDifferentCurrencies.Error()})
		return
	}

	if originAcc.Balance.Amount.LessThan(in.Amount) {
		c.JSON(400, gin.H{"error": errs.ErrAccountInsufficientBalance.Error()})
		return
	}

	// Withdrawing money from the origin account
	withdrawTrx, err := ctrl.AccountsService.Withdraw(c, dto.WithdrawCommand{
		Account: originAcc,
		Money: *types.NewMoney(in.Amount, originAcc.Balance.Currency),
		Comment: fmt.Sprintf("Перевод на %v", in.DestinationAccountID),
	})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// TODO: While deposit could fail, withdraw is already done
	// Find a way to handle withdraw cancelling (DB transactions)

	// Depositing money to the destination account
	depositTrx, err := ctrl.AccountsService.Deposit(c, dto.DepositCommand{
		Account: destAcc,
		Money: *types.NewMoney(in.Amount, destAcc.Balance.Currency),
		Comment: fmt.Sprintf("Перевод от %v", in.OriginAccountID),
	})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"deposit_transaction": depositTrx, "withdraw_transaction": withdrawTrx})
	return
}
