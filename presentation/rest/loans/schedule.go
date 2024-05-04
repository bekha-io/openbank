package loans

import (
	"net/http"
	"time"

	"github.com/bekha-io/openbank/domain/dto"
	"github.com/bekha-io/openbank/domain/types"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func (ctrl *LoanController) AnnuitySchedule(c *gin.Context) {
	type req struct {
		Amount        float64        `json:"amount" binding:"required"`
		Currency      types.Currency `json:"currency" binding:"required"`
		Duration      uint           `json:"duration" binding:"required"`
		InterestRate  float64        `json:"interest_rate" binding:"required"`
		RepayStartsAt time.Time      `json:"repay_starts_at" binding:"required"`
	}

	var in req
	err := c.BindJSON(&in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	installments, err := ctrl.LoanService.CalculateAnnuityInstallments(c, dto.CalculateAnnuityInstallmentsQuery{
		LoanAmount:    *types.NewMoney(decimal.NewFromFloat(in.Amount), in.Currency),
		InterestRate:  decimal.NewFromFloat(in.InterestRate),
		RepayStartsAt: in.RepayStartsAt,
		Duration:      in.Duration,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, installments)
}
