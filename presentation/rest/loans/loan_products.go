package loans

import (
	"net/http"

	"github.com/bekha-io/openbank/domain/dto"
	"github.com/bekha-io/openbank/domain/types"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func (ctrl *LoanController) GetLoanProducts(c *gin.Context) {
	products, err := ctrl.LoanService.GetLoanProducts(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, products)
}

func (ctrl *LoanController) CreateLoanProduct(c *gin.Context) {
	type req struct {
		Name                string          `json:"name"`
		MinDuration         uint            `json:"min_duration"`
		MaxDuration         uint            `json:"max_duration"`
		Currency            types.Currency  `json:"currency"`
		MinAmount           decimal.Decimal `json:"min_amount"`
		MaxAmount           decimal.Decimal `json:"max_amount"`
		InterestRate        decimal.Decimal `json:"interest_rate"`
		LoanType            types.LoanType  `json:"loan_type"`
		DailyOverduePenalty decimal.Decimal `json:"daily_overdue_penalty"`
	}

	var in req
	err := c.BindJSON(&in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := dto.CreateLoanProductCommand{
		Name:                in.Name,
		MinDuration:         in.MinDuration,
		MaxDuration:         in.MaxDuration,
		MinAmount:           types.Money{Amount: in.MinAmount, Currency: in.Currency},
		MaxAmount:           types.Money{Amount: in.MaxAmount, Currency: in.Currency},
		InterestRate:        in.InterestRate,
		LoanType:            in.LoanType,
		DailyOverduePenalty: in.DailyOverduePenalty,
	}
	err = ctrl.LoanService.CreateLoanProduct(c, cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.String(201, "loan product created")
}
