package services

import (
	"context"

	"github.com/bekha-io/openbank/domain/dto"
	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/repository"
	"github.com/bekha-io/openbank/domain/types"
	"github.com/shopspring/decimal"
)

type ILoanService interface {
	CreateLoanProduct(ctx context.Context, cmd dto.CreateLoanProductCommand) error
	GetLoanProducts(ctx context.Context) ([]*entities.LoanProduct, error)
	CalculateAnnuityInstallments(ctx context.Context, q dto.CalculateAnnuityInstallmentsQuery) ([]entities.LoanInstallment, error)
}

var _ ILoanService = (*LoanService)(nil)

type LoanService struct {
	LoanRepository repository.ILoanRepository
}

func NewLoanService(lr repository.ILoanRepository) *LoanService {
	return &LoanService{
		LoanRepository: lr,
	}
}

// GetLoanProducts implements ILoanService.
func (s *LoanService) GetLoanProducts(ctx context.Context) ([]*entities.LoanProduct, error) {
	return s.LoanRepository.GetAllLoanProducts(ctx)
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

// CreateLoanProduct implements ILoanService.
func (s *LoanService) CreateLoanProduct(ctx context.Context, cmd dto.CreateLoanProductCommand) error {
	// Validate input parameters
	if err := cmd.Validate(); err != nil {
		return err
	}

	lp, err := entities.NewLoanProduct(cmd.Name, cmd.MinAmount.Amount, cmd.MaxAmount.Amount,
		cmd.MinAmount.Currency, cmd.InterestRate, cmd.LoanType, cmd.DailyOverduePenalty, cmd.MinDuration, cmd.MaxDuration)
	if err != nil {
		return err
	}

	err = s.LoanRepository.SaveLoanProduct(ctx, lp)
	if err != nil {
		return err
	}
	return nil
}
