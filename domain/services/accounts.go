package services

import (
	"context"
	"errors"

	"github.com/bekha-io/vaultonomy/domain/dto"
	"github.com/bekha-io/vaultonomy/domain/entities"
	"github.com/bekha-io/vaultonomy/domain/repository"
	"github.com/bekha-io/vaultonomy/domain/types"
	"github.com/bekha-io/vaultonomy/domain/types/errs"
)

type IAccountService interface {
	CreateAccount(ctx context.Context, cmd dto.CreateAccountCommand) error
	GetAccountByID(ctx context.Context, id types.AccountID) (*entities.Account, error)
	Deposit(ctx context.Context, cmd dto.DepositCommand) (*entities.Transaction, error)
	Withdraw(ctx context.Context, cmd dto.WithdrawCommand) (*entities.Transaction, error)
}

var _ IAccountService = (*AccountsService)(nil)

type AccountsService struct {
	AccountsRepo     repository.IAccountRepository
	TransactionsRepo repository.ITransactionRepository
}

func NewAccountsService(accountsRepo repository.IAccountRepository, transactionsRepo repository.ITransactionRepository) *AccountsService {
	return &AccountsService{
		AccountsRepo:     accountsRepo,
		TransactionsRepo: transactionsRepo,
	}
}

func (s *AccountsService) CreateAccount(ctx context.Context, cmd dto.CreateAccountCommand) error {
	acc := entities.NewAccount(cmd.CustomerID, cmd.Currency)
	err := s.AccountsRepo.Save(ctx, acc)
	if err != nil {
		return err
	}
	return nil
}

func (s *AccountsService) GetAccountByID(ctx context.Context, id types.AccountID) (*entities.Account, error) {
	acc, err := s.AccountsRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Join(errs.ErrAccountNotFound, err)
	}
	return acc, err
}

// Deposit implements IAccountService.
func (s *AccountsService) Deposit(ctx context.Context, cmd dto.DepositCommand) (*entities.Transaction, error) {
	err := cmd.Money.Validate()
	if err != nil {
		return nil, err
	}

	if cmd.Money.Currency != cmd.Account.Balance.Currency {
		return nil, errs.ErrAccountDifferentCurrencies
	}

	err = cmd.Account.Balance.Add(cmd.Money)
	if err != nil {
		return nil, err
	}

	err = s.AccountsRepo.Save(ctx, cmd.Account)
	if err != nil {
		return nil, err
	}

	transaction := entities.NewTransaction(cmd.Account.ID, types.DepositTransactionType, cmd.Money)
	transaction.Comment = cmd.Comment
	err = s.TransactionsRepo.Save(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// Withdraw implements IAccountService.
func (s *AccountsService) Withdraw(ctx context.Context, cmd dto.WithdrawCommand) (*entities.Transaction, error) {
	err := cmd.Money.Validate()
	if err != nil {
		return nil, err
	}

	if cmd.Money.Currency != cmd.Account.Balance.Currency {
		return nil, errs.ErrAccountDifferentCurrencies
	}

	err = cmd.Account.Balance.Sub(cmd.Money)
	if err != nil {
		return nil, err
	}

	err = s.AccountsRepo.Save(ctx, cmd.Account)
	if err != nil {
		return nil, err
	}

	transaction := entities.NewTransaction(cmd.Account.ID, types.WithdrawTransactionType, cmd.Money)
	transaction.Comment = cmd.Comment
	err = s.TransactionsRepo.Save(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
