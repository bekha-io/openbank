package fineract

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *FineractClient) CreateClient(ctx context.Context, in CreateClientIn) (CreateClientOut, error) {
	body, err := json.Marshal(in)
	if err != nil {
		return CreateClientOut{}, err
	}

	out := &CreateClientOut{}
	if err := c.SendRequest(ctx, "POST", "/v1/clients", body, &out); err != nil {
		return CreateClientOut{}, fmt.Errorf("cannot create client: %v", err)
	}

	return *out, nil
}

func (c *FineractClient) GetClientsAccounts(ctx context.Context, clientId uint) (GetClientsAccountsOut, error) {
	out := &GetClientsAccountsOut{}
	if err := c.SendRequest(ctx, "GET", fmt.Sprintf("/v1/clients/%d/accounts", clientId), nil, &out); err != nil {
		return GetClientsAccountsOut{}, fmt.Errorf("cannot get clients accounts: %v", err)
	}
	return *out, nil
}

func (c *FineractClient) GetClientById(ctx context.Context, id uint) (Client, error) {
	out := &Client{}
	if err := c.SendRequest(ctx, "GET", fmt.Sprintf("/v1/clients/%d", id), nil, &out); err != nil {
		return Client{}, fmt.Errorf("cannot get client by id: %v", err)
	}
	return *out, nil
}

func (c *FineractClient) GetClientByExternalId(ctx context.Context, clientId string) (Client, error) {
	out := &Client{}
	if err := c.SendRequest(ctx, "GET", fmt.Sprintf("/v1/clients/external-id/%s", clientId), nil, &out); err != nil {
		return Client{}, fmt.Errorf("cannot get client by external id: %v", err)
	}
	return *out, nil
}
