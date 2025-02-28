package services

import (
	"context"
	"errors"
	"time"

	"github.com/bekha-io/openbank/domain/dto"
	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/repository"
	"github.com/bekha-io/openbank/domain/types/errs"
	"github.com/bekha-io/openbank/infrastructure/fineract"
	"github.com/shopspring/decimal"
)

type TransferIn struct {
	FromAccountId   uint
	ToBeneficiaryId uint
	Amount          decimal.Decimal
	Comment         string
}

type IAccountService interface {
	CreateAccount(ctx context.Context, cmd dto.CreateAccountCommand) error
	GetAccountByID(ctx context.Context, id uint) (*entities.Account, error)
	Transfer(ctx context.Context, in TransferIn) (*entities.Transaction, error)
	GetAccountTransactions(ctx context.Context, id uint) ([]*entities.Transaction, error)
}

var _ IAccountService = (*AccountsService)(nil)

type AccountsService struct {
	Core             fineract.FineractClientI
	BeneficiaryRepo  repository.IBenificiaryRepository
	TransactionsRepo repository.ITransactionRepository
}

// Transfer implements IAccountService.
func (s *AccountsService) Transfer(ctx context.Context, in TransferIn) (*entities.Transaction, error) {

	beneficiary, err := s.BeneficiaryRepo.GetBeneficiaryByID(ctx, in.ToBeneficiaryId)
	if err != nil {
		return nil, err
	}

	accounts, err := s.Core.GetClientsAccounts(ctx, beneficiary.BeneficiaryCustomerID)
	if err != nil {
		return nil, err
	}

	if len(accounts.SavingsAccounts) == 0 {
		return nil, errs.ErrAccountNotFound
	}

	var primaryAccountId uint = accounts.SavingsAccounts[0].ID

	out, err := s.Core.Transfer(ctx, fineract.TransferIn{
		FromOfficeId:        1,
		ToOfficeId:          1,
		FromAccountType:     2,
		ToAccountType:       2,
		FromClientId:        beneficiary.OwnerCustomerID,
		ToClientId:          beneficiary.BeneficiaryCustomerID,
		FromAccountId:       in.FromAccountId,
		ToAccountId:         primaryAccountId,
		TransferAmount:      in.Amount.InexactFloat64(),
		DateFormat:          "dd.MM.yyyy HH:mm:ss",
		TransferDate:        time.Now().Format("02.01.2006 15:04:05"),
		TransferDescription: "Перевод между счетами",
		Locale:              "ru",
	})

	tr := &entities.Transaction{
		ID:            out.ResourceId,
		Amount:        in.Amount,
		FromAccountId: in.FromAccountId,
		ToAccountId:   primaryAccountId,
		CreatedAt:     time.Now().UTC(),
		Comment:       in.Comment,
	}
	if err != nil {
		tr.Status = entities.TransactionStatusFailed
		tr.StatusReason = err.Error()
	}

	err = s.TransactionsRepo.SaveTransaction(ctx, tr)
	return tr, err
}

func NewAccountsService(core fineract.FineractClientI, br repository.IBenificiaryRepository, tr repository.ITransactionRepository) *AccountsService {
	return &AccountsService{
		Core:             core,
		BeneficiaryRepo:  br,
		TransactionsRepo: tr,
	}
}

func (s *AccountsService) CreateAccount(ctx context.Context, cmd dto.CreateAccountCommand) error {
	_, err := s.Core.CreateSavingsAccount(ctx, fineract.CreateSavingsAccountIn{
		ClientId:        uint(cmd.CustomerID),
		DateFormat:      "dd.MM.yyyy",
		Locale:          "ru",
		ProductId:       1,
		SubmittedOnDate: time.Now().Format("02.01.2006"),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *AccountsService) GetAccountByID(ctx context.Context, id uint) (*entities.Account, error) {
	acc, err := s.Core.GetSavingsAccountById(ctx, id)
	if err != nil {
		return nil, errors.Join(errs.ErrAccountNotFound, err)
	}
	return acc.ToEntity(), err
}

func (s *AccountsService) GetAccountTransactions(ctx context.Context, id uint) ([]*entities.Transaction, error) {
	transactions, err := s.TransactionsRepo.GetTransactionsByAccountID(ctx, id)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
