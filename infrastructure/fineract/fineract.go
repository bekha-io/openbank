package fineract

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/types"
	"github.com/shopspring/decimal"
)

type Non200Response struct {
	Errors []struct {
		DeveloperMessage             string `json:"developerMessage"`
		DefaultUserMessage           string `json:"defaultUserMessage"`
		UserMessageGlobalisationCode string `json:"userMessageGlobalisationCode"`
	} `json:"errors"`
}

type Currency struct {
	Code          string `json:"code"`
	DecimalPlaces uint   `json:"decimalPlaces"`
	DisplayLabel  string `json:"displayLabel"`
	Name          string `json:"name"`
	NameCode      string `json:"nameCode"`
}

type SavingsAccountTransaction struct {
	AccountId uint     `json:"accountId"`
	AccountNo string   `json:"accountNo"`
	Amount    float64  `json:"amount"`
	Currency  Currency `json:"currency"`
	Date      string   `json:"date"`
	EntryType string   `json:"entryType"`
	Id        uint     `json:"id"`
}

type GetSavingsAccountTransactionsOut struct {
	Content []SavingsAccountTransaction `json:"content"`
	Total   uint                        `json:"total"`
}

type Client struct {
	AccountNo          string    `json:"accountNo"`
	ActivationDate     TimeSlice `json:"activationDate"`
	DateOfBirth        TimeSlice `json:"dateOfBirth"`
	Active             bool      `json:"active"`
	DisplayName        string    `json:"displayName"`
	EmailAddress       string    `json:"emailAddress"`
	ExternalId         string    `json:"externalId"`
	Firstname          string    `json:"firstname"`
	Lastname           string    `json:"lastname"`
	Id                 uint      `json:"id"`
	OfficeId           uint      `json:"officeId"`
	OfficeName         string    `json:"officeName"`
	SavingsProductId   uint      `json:"savingsProductId"`
	SavingsProductName string    `json:"savingsProductName"`
	Status             struct {
		Code        string `json:"code"`
		Description string `json:"description"`
		Id          uint   `json:"id"`
	} `json:"status"`
}

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
	ActivationDate   string `json:"activationDate,omitempty"`
	Active           bool   `json:"active,omitempty"`
	Locale           string `json:"locale,omitempty"`
	MobileNo         string `json:"mobileNo,omitempty"`
	OfficeId         uint   `json:"officeId,omitempty"`
	LastName         string `json:"lastname,omitempty"`
	MiddleName       string `json:"middlename,omitempty"`
	FirstName        string `json:"firstname,omitempty"`
	ExternalId       string `json:"externalId,omitempty"`
	DateFormat       string `json:"dateFormat,omitempty"`
	SavingsProductId uint   `json:"savingsProductId,omitempty"`
}

/*
Submits new savings application

Mandatory Fields: clientId or groupId, productId, submittedOnDate

Optional Fields: accountNo, externalId, fieldOfficerId

Inherited from Product (if not provided): nominalAnnualInterestRate, interestCompoundingPeriodType, interestCalculationType, interestCalculationDaysInYearType, minRequiredOpeningBalance, lockinPeriodFrequency, lockinPeriodFrequencyType, withdrawalFeeForTransfers, allowOverdraft, overdraftLimit, withHoldTax

Additional Mandatory Field if Entity-Datatable Check is enabled for the entity of type Savings: datatables
*/
type CreateSavingsAccountIn struct {
	ClientId        uint   `json:"clientId"`
	DateFormat      string `json:"dateFormat"`
	ExternalId      string `json:"externalId,omitempty"`
	Locale          string `json:"locale,omitempty"`
	ProductId       uint   `json:"productId"`
	SubmittedOnDate string `json:"submittedOnDate"`
}

type CreateSavingsAccountOut struct {
	ClientId   uint `json:"clientId"`
	OfficeId   uint `json:"officeId"`
	ResourceId uint `json:"resourceId"`
	SavingsId  uint `json:"savingsId"`
}

type GetClientsAccountsOut struct {
	SavingsAccounts []struct {
		ID          uint   `json:"id"`
		AccountNo   string `json:"accountNo,omitempty"`
		ProductId   uint   `json:"productId,omitempty"`
		ProductName string `json:"productName,omitempty"`
		Currency    struct {
			Code          string `json:"code"`
			Name          string `json:"name"`
			DecimalPlaces uint   `json:"decimalPlaces"`
		}
	} `json:"savingsAccounts"`
}

type SavingsAccount struct {
	AccountNo          string   `json:"accountNo,omitempty"`
	ClientId           uint     `json:"clientId,omitempty"`
	ClientName         string   `json:"clientName,omitempty"`
	Currency           Currency `json:"currency,omitempty"`
	FieldOfficerId     uint     `json:"fieldOfficerId"`
	Id                 uint     `json:"id"`
	SavingsProductId   uint     `json:"savingsProductId"`
	SavingsProductName string   `json:"savingsProductName"`
	Status             struct {
		Active   bool `json:"active"`
		Approved bool `json:"approved"`
		Closed   bool `json:"closed"`
		Rejected bool `json:"rejected"`
	} `json:"status"`
	Summary struct {
		AccountBalance   float64  `json:"accountBalance"`
		AvailableBalance float64  `json:"availableBalance"`
		Currency         Currency `json:"currency"`
	}
}

func (s SavingsAccount) ToEntity() *entities.Account {
	return &entities.Account{
		ID:         s.Id,
		AccountNo:  s.AccountNo,
		CustomerID: s.ClientId,
		Balance:    types.NewMoney(decimal.NewFromFloat(s.Summary.AccountBalance), types.Currency(s.Summary.Currency.Code)),
	}
}

/*
{
  "dateFormat": "dd MMMM yyyy",
  "fromAccountId": 1,
  "fromAccountType": 2,
  "fromClientId": 1,
  "fromOfficeId": 1,
  "locale": "en",
  "toAccountId": 2,
  "toAccountType": 2,
  "toClientId": 1,
  "toOfficeId": 1,
  "transferAmount": 112.45,
  "transferDate": "01 August 2011",
  "transferDescription": "A description of the transfer"
}
*/
// TransferIn структура для запроса на перевод
type TransferIn struct {
	FromAccountId   uint `json:"fromAccountId"`
	ToAccountId     uint `json:"toAccountId"`
	FromAccountType uint `json:"fromAccountType"`
	ToAccountType   uint `json:"toAccountType"`
	FromOfficeId    uint `json:"fromOfficeId"`
	ToOfficeId      uint `json:"toOfficeId"`
	FromClientId    uint `json:"fromClientId"`
	ToClientId      uint `json:"toClientId"`

	TransferAmount      float64 `json:"transferAmount"`
	DateFormat          string  `json:"dateFormat"`
	Locale              string  `json:"locale"`
	TransferDate        string  `json:"transferDate"`
	TransferDescription string  `json:"transferDescription"`
}

type TransferOut struct {
	SavingsId  uint `json:"savingsId"`
	ResourceId uint `json:"resourceId"`
}

type FineractClientI interface {
	// Clients
	CreateClient(ctx context.Context, in CreateClientIn) (CreateClientOut, error)
	GetClientsAccounts(ctx context.Context, clientId uint) (GetClientsAccountsOut, error)
	GetClientByExternalId(ctx context.Context, clientId string) (Client, error)
	GetClientById(ctx context.Context, id uint) (Client, error)

	// Savings accounts
	CreateSavingsAccount(ctx context.Context, in CreateSavingsAccountIn) (CreateSavingsAccountOut, error)
	GetSavingsAccountById(ctx context.Context, id uint) (SavingsAccount, error)
	GetSavingsAccountTransactions(ctx context.Context, id uint) (GetSavingsAccountTransactionsOut, error)
	Transfer(ctx context.Context, in TransferIn) (TransferOut, error)
}

type FineractClient struct {
	BaseUrl    string
	TenantName string
	Username   string
	Password   string
}

func (c *FineractClient) SendRequest(ctx context.Context, method string, path string, reqBody interface{}, respTarget interface{}, reqModifier ...func(r *http.Request)) error {
	cl := &http.Client{}

	// Preparing request body
	bodyJson, err := json.Marshal(&reqBody)
	if err != nil {
		return fmt.Errorf("cannot marshal request body to JSON: %v", err)
	}

	// Creating an HTTP request
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%v%v", c.BaseUrl, path), bytes.NewReader(bodyJson))
	if err != nil {
		return fmt.Errorf("cannot create HTTP request: %v", err)
	}
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Fineract-Platform-TenantId", c.TenantName)

	for _, rm := range reqModifier {
		rm(req)
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

	if resp.StatusCode >= 400 {
		// Try to parse errors
		var fe Non200Response
		err = json.Unmarshal(bodyBytes, &fe)
		if err == nil {
			return fmt.Errorf("HTTP request failed with status code %d: %s (errors: %v)", resp.StatusCode, string(bodyBytes), fe.Errors)
		}

		var err error
		for _, e := range fe.Errors {
			err = errors.Join(err, handleErrorCode(e.UserMessageGlobalisationCode))
		}
		return err
	}

	return nil
}
