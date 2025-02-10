package fineract

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client successfully created response
type CreateClientOut struct {
	ClientId           uint `json:"clientId"`
	GroupId            uint `json:"groupId"`
	OfficeId           uint `json:"officeId"`
	ResourceExternalId uint `json:"resourceExternalId"`
	ResourceID         uint `json:"resourceId"`
}

/*
Note:
You can enter either:firstname/middlename/lastname - for a person (middlename is optional) OR fullname - for a business or organisation (or person known by one name).
2.If address is enable(enable-address=true), then additional field called address has to be passed.
Mandatory Fields: firstname and lastname OR fullname, officeId, active=true and activationDate OR active=false, if(address enabled) address
Optional Fields: groupId, externalId, accountNo, staffId, mobileNo, savingsProductId, genderId, clientTypeId, clientClassificationId
*/
type CreateClientIn struct {
	ActivationDate string `json:"activationDate,omitempty"`
	Active         bool   `json:"active,omitempty"`
	Locale         string `json:"locale,omitempty"`
	MobileNo       string `json:"mobileNo,omitempty"`
	OfficeId       string `json:"officeId,omitempty"`
	LastName       string `json:"lastname,omitempty"`
	MiddleName     string `json:"middlename,omitempty"`
	FirstName      string `json:"firstname,omitempty"`
	ExternalId     string `json:"externalId,omitempty"`
	DateFormat     string `json:"dateFormat,omitempty"`
}

type CreateSavingsAccountIn struct {
}

type CreateSavingsAccountOut struct {
}

type GetClientsAccountsOut struct {
	SavingsAccounts []struct{
		ID uint `json:"id"`
		AccountNo string `json:"accountNo,omitempty"`
		ProductId uint `json:"productId,omitempty"`
		ProductName string `json:"productName,omitempty"`
		Currency struct {
			Code string `json:"code"`
            Name string `json:"name"`
			DecimalPlaces uint `json:"decimalPlaces"`
			
		}
	} `json:"savingsAccounts"`
}

type FineractClientI interface {
	CreateClient(ctx context.Context, in CreateClientIn) (*CreateClientOut, error)
	CreateSavingsAccount(ctx context.Context, in CreateSavingsAccountIn)
	GetClientsAccounts(ctx context.Context, clientId string) (GetClientsAccountsOut, error)
}

type FineractClient struct {
	BaseUrl    string
	TenantName string
	Username   string
	Password   string
}

func (c *FineractClient) SendRequest(ctx context.Context, method string, path string, reqBody interface{}, respTarget interface{}) error {
	cl := &http.Client{}

	// Preparing request body
	bodyJson, err := json.Marshal(&reqBody)
	if err != nil {
		return fmt.Errorf("cannot marshal request body to JSON: %v", err)
	}

	// Creating an HTTP request
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, method, path, bytes.NewReader(bodyJson))
	if err != nil {
		return fmt.Errorf("cannot create HTTP request: %v", err)
	}

	resp, err := cl.Do(req)
	if err != nil {
		return fmt.Errorf("cannot send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Decoding response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %v", err)
	}

	// Unmarshalling response body to target struct
	err = json.Unmarshal(bodyBytes, respTarget)
	if err != nil {
		return fmt.Errorf("cannot unmarshal response body to target struct: %v | response body: %v", err, string(bodyBytes))
	}

	return nil
}
