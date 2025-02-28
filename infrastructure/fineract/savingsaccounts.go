package fineract

import "fmt"

import "context"

func (f *FineractClient) CreateSavingsAccount(ctx context.Context, in CreateSavingsAccountIn) (CreateSavingsAccountOut, error) {
	out := &CreateSavingsAccountOut{}
	if err := f.SendRequest(ctx, "POST", "/v1/savingsaccounts", in, &out); err != nil {
		return CreateSavingsAccountOut{}, err
	}
	return *out, nil
}

func (f *FineractClient) GetSavingsAccountById(ctx context.Context, id uint) (SavingsAccount, error) {
	out := &SavingsAccount{}
	if err := f.SendRequest(ctx, "GET", "/v1/savingsaccounts/"+fmt.Sprint(id), nil, &out); err != nil {
		return SavingsAccount{}, err
	}
	return *out, nil
}

func (f *FineractClient) GetSavingsAccountTransactions(ctx context.Context, id uint) (GetSavingsAccountTransactionsOut, error) {
	out := &GetSavingsAccountTransactionsOut{}
	if err := f.SendRequest(ctx, "GET", "/v1/savingsaccounts/"+fmt.Sprint(id)+"/transactions/search", nil, &out); err != nil {
		return GetSavingsAccountTransactionsOut{}, err
	}
	return *out, nil
}

func (f *FineractClient) Transfer(ctx context.Context, in TransferIn) (TransferOut, error) {
	out := &TransferOut{}
	if err := f.SendRequest(ctx, "POST", "/v1/accounttransfers", in, &out); err != nil {
		return TransferOut{}, err
	}
	return TransferOut{}, nil
}
