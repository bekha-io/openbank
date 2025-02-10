package fineract

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *FineractClient) CreateClient(ctx context.Context, in CreateClientIn) (*CreateClientOut, error) {
	body, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	out := &CreateClientOut{}
	if err := c.SendRequest(ctx, "POST", "/v1/clients", body, &out); err != nil {
		return nil, fmt.Errorf("cannot create client: %v", err)
	}

	return out, nil
}
