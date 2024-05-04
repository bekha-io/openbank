package entities

import (
	"time"

	"github.com/bekha-io/openbank/domain/types"
	"github.com/shopspring/decimal"
)

type Loan struct {
	LoanID        types.LoanID
	LoanProductID types.LoanProductID
	CustomerID    types.CustomerID

	// Amount of money provided as loan to the customer
	Amount types.Money

	// Duration specifies the amount of days the loan dues
	Duration uint // (in days)

	// RepayStartsAt specifies the date and time of the first installment to be repaid
	RepayStartsAt time.Time

	// RepayEndsAt specifies the date and time of the last installment to be repaid
	RepayEndsAt time.Time

	// Timestamp of the loan record created within the system
	CreatedAt time.Time

	// Installments specifies the scheduled installments of the given loan
	Installments []LoanInstallment
}

type LoanInstallment struct {
	// RepayAmount specifies the amount of money to be paid within the given installment
	RepayAmount types.Money `json:"repay_amount"`

	// Principal specifies the portion of total repay amount being paid as the loan principal
	Principal decimal.Decimal `json:"principal"`

	// Interest specifies the portion of total repay amount being paid as the loan interest
	Interest decimal.Decimal `json:"interest"`

	// RepayAmountAfter specifies the amount of money left to be paid paid after the given installment fully repaid
	RepayAmountAfter types.Money `json:"repay_amount_after"`

	// RepayAt specifies the date installment should be repaid
	RepayAt  time.Time `json:"repay_at"`
	IsRepaid bool      `json:"is_repaid"`

	// Payments shows detailed information about repays within the given installment
	Payments []LoanInstallmentPayment `json:"payments"`
}

func (i LoanInstallment) AmountLeft() *types.Money {
	var paidAmount decimal.Decimal

	if len(i.Payments) == 0 {
		return &i.RepayAmount
	}

	for _, p := range i.Payments {
		paidAmount = paidAmount.Add(p.Amount.Amount)
	}

	return &types.Money{
		Amount:   i.RepayAmount.Amount.Sub(paidAmount),
		Currency: i.RepayAmount.Currency,
	}
}

type LoanInstallmentPayment struct {
	Amount      *types.Money
	PaymentType types.PaymentType
	PaidAt      time.Time
}
