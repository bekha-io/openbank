package services

import (
	"context"

	"github.com/bekha-io/vaultonomy/domain/dto"
	"github.com/bekha-io/vaultonomy/domain/entities"
	"github.com/bekha-io/vaultonomy/domain/types"
	"github.com/shopspring/decimal"
)

type ILoanService interface {
	CalculateAnnuityInstallments(ctx context.Context, q dto.CalculateAnnuityInstallmentsQuery) ([]entities.LoanInstallment, error)
}

var _ ILoanService = (*LoanService)(nil)

type LoanService struct{}

func NewLoanService() *LoanService {
	return &LoanService{}
}

func (s *LoanService) CalculateAnnuityInstallments(ctx context.Context, q dto.CalculateAnnuityInstallmentsQuery) ([]entities.LoanInstallment, error) {
	// Рассчитываем ежемесячную процентную ставку
	monthlyInterestRate := q.InterestRate.Div(decimal.NewFromInt(12))

	// Рассчитываем аннуитетный коэффициент
	annuityCoefficient := monthlyInterestRate.Mul(q.LoanAmount.Amount).Div(decimal.NewFromInt(1).Sub(decimal.NewFromInt(1).Add(monthlyInterestRate).Pow(decimal.NewFromInt(-int64(q.Duration)))))

	// Создаем массив для хранения платежей
	schedule := make([]entities.LoanInstallment, q.Duration)

	// Инициализируем остаток основного долга
	remainingPrincipal := q.LoanAmount.Amount

	repayDate := q.RepayStartsAt

	// Рассчитываем каждый платеж
	for i := 0; i < int(q.Duration); i++ {
		// Рассчитываем сумму платежа
		paymentAmount := annuityCoefficient

		// Рассчитываем процентную часть платежа
		interestPayment := remainingPrincipal.Mul(monthlyInterestRate)

		// Рассчитываем основную часть платежа
		principalPayment := paymentAmount.Sub(interestPayment)

		// Обновляем остаток основного долга
		remainingPrincipal = remainingPrincipal.Sub(principalPayment)

		repayDate = q.RepayStartsAt.AddDate(0, i, 0)

		// Создаем новый платеж
		installment := entities.LoanInstallment{
			RepayAmount:      types.Money{Amount: paymentAmount.Round(2), Currency: q.LoanAmount.Currency},
			Principal:        principalPayment.Round(2),
			Interest:         interestPayment.Round(2),
			RepayAmountAfter: types.Money{Amount: remainingPrincipal.Round(2), Currency: q.LoanAmount.Currency},
			RepayAt:          repayDate.Round(2),
			IsRepaid:         false,
		}

		// Добавляем платеж в график
		schedule[i] = installment
	}

	return schedule, nil
}
