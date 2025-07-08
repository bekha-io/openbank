package me

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bekha-io/openbank/domain/dto"
	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/repository"
	"github.com/bekha-io/openbank/domain/services"
	"github.com/bekha-io/openbank/domain/types"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Упрощенные моки без testify/mock
type MockAccountsService struct {
	GetAccountByIDFunc         func(ctx context.Context, id uint) (*entities.Account, error)
	GetAccountTransactionsFunc func(ctx context.Context, in repository.GetTransactionsIn) ([]*entities.Transaction, error)
	CreateAccountFunc          func(ctx context.Context, cmd dto.CreateAccountCommand) error
	TransferFunc               func(ctx context.Context, in services.TransferIn) (*entities.Transaction, error)
}

func (m *MockAccountsService) GetAccountByID(ctx context.Context, id uint) (*entities.Account, error) {
	if m.GetAccountByIDFunc != nil {
		return m.GetAccountByIDFunc(ctx, id)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *MockAccountsService) GetAccountTransactions(ctx context.Context, in repository.GetTransactionsIn) ([]*entities.Transaction, error) {
	if m.GetAccountTransactionsFunc != nil {
		return m.GetAccountTransactionsFunc(ctx, in)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *MockAccountsService) CreateAccount(ctx context.Context, cmd dto.CreateAccountCommand) error {
	if m.CreateAccountFunc != nil {
		return m.CreateAccountFunc(ctx, cmd)
	}
	return fmt.Errorf("not implemented")
}

func (m *MockAccountsService) Transfer(ctx context.Context, in services.TransferIn) (*entities.Transaction, error) {
	if m.TransferFunc != nil {
		return m.TransferFunc(ctx, in)
	}
	return nil, fmt.Errorf("not implemented")
}

type MockCustomersService struct {
	GetCustomerFunc         func(ctx context.Context, id uint) (*entities.Customer, error)
	GetCustomerAccountsFunc func(ctx context.Context, customer entities.Customer) ([]*entities.Account, error)
}

func (m *MockCustomersService) CreateCustomer(ctx context.Context, in dto.CreateIndividualCustomerCommand) error {
	return fmt.Errorf("not implemented")
}

func (m *MockCustomersService) GetCustomer(ctx context.Context, id uint) (*entities.Customer, error) {
	if m.GetCustomerFunc != nil {
		return m.GetCustomerFunc(ctx, id)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *MockCustomersService) GetCustomerByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.Customer, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MockCustomersService) GetCustomerAccounts(ctx context.Context, customer entities.Customer) ([]*entities.Account, error) {
	if m.GetCustomerAccountsFunc != nil {
		return m.GetCustomerAccountsFunc(ctx, customer)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *MockCustomersService) GetCustomerBeneficiaries(ctx context.Context, id uint) ([]*entities.Benificiary, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MockCustomersService) GetBeneficiaryByID(ctx context.Context, id uint) (*entities.Benificiary, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MockCustomersService) CreateBeneficiary(ctx context.Context, in services.CreateBenificiaryIn) (*entities.Benificiary, error) {
	return nil, fmt.Errorf("not implemented")
}

func setupTestController() (*Controller, *MockAccountsService, *MockCustomersService) {
	mockAccountsService := &MockAccountsService{}
	mockCustomersService := &MockCustomersService{}

	controller := &Controller{
		AccountsService:  mockAccountsService,
		CustomersService: mockCustomersService,
	}

	return controller, mockAccountsService, mockCustomersService
}

func setupGinTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// Настройка реального HTTP сервера для тестирования
func setupHTTPTestServer() (*httptest.Server, *Controller) {
	gin.SetMode(gin.TestMode)

	mockAccountsService := &MockAccountsService{}
	mockCustomersService := &MockCustomersService{}

	controller := &Controller{
		AccountsService:  mockAccountsService,
		CustomersService: mockCustomersService,
	}

	router := gin.New()

	// Middleware для установки customer ID
	router.Use(func(c *gin.Context) {
		c.Set("customerId", float64(456))
		c.Next()
	})

	router.GET("/accounts/:id/transactions", controller.GetAccountTransactions)

	server := httptest.NewServer(router)
	return server, controller
}

func TestController_GetAccountTransactions_WithHTTPServer(t *testing.T) {
	server, controller := setupHTTPTestServer()
	defer server.Close()

	// Тестовые данные
	testAccount := &entities.Account{
		ID:         123,
		CustomerID: 456,
		Balance:    types.NewMoney(decimal.NewFromFloat(1000.0), types.Currency("UZS")),
	}

	testTransactions := []*entities.Transaction{
		{
			ID:            1,
			ExternalID:    "ext-1",
			FromAccountId: 123,
			ToAccountId:   456,
			Category:      types.TransactionCategoryTransfer,
			Status:        entities.TransactionStatusCompleted,
			Amount:        decimal.NewFromFloat(100.50),
			Comment:       "Test transaction 1",
			CreatedAt:     time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		},
		{
			ID:            2,
			ExternalID:    "ext-2",
			FromAccountId: 456,
			ToAccountId:   123,
			Category:      types.TransactionCategoryTransfer,
			Status:        entities.TransactionStatusCompleted,
			Amount:        decimal.NewFromFloat(200.75),
			Comment:       "Test transaction 2",
			CreatedAt:     time.Date(2024, 1, 20, 15, 45, 0, 0, time.UTC),
		},
	}

	tests := []struct {
		name           string
		url            string
		mockSetup      func()
		expectedStatus int
		expectedCount  int
		shouldHaveErr  bool
	}{
		{
			name: "успешное получение транзакций без фильтров",
			url:  server.URL + "/accounts/123/transactions",
			mockSetup: func() {
				// Мокируем внутренние сервисы
				controller.AccountsService.(*MockAccountsService).GetAccountByIDFunc = func(ctx context.Context, id uint) (*entities.Account, error) {
					if id == 123 {
						return testAccount, nil
					}
					return nil, fmt.Errorf("account not found")
				}
				controller.AccountsService.(*MockAccountsService).GetAccountTransactionsFunc = func(ctx context.Context, in repository.GetTransactionsIn) ([]*entities.Transaction, error) {
					if in.AccountID == 123 {
						return testTransactions, nil
					}
					return nil, fmt.Errorf("unexpected account ID: %d", in.AccountID)
				}
			},
			expectedStatus: 200,
			expectedCount:  2,
		},
		{
			name: "фильтрация по дате с timezone",
			url:  server.URL + "/accounts/123/transactions?date_from=2024-01-01T10:00:00Z&date_to=2024-01-31T23:59:59%2B05:00",
			mockSetup: func() {
				controller.AccountsService.(*MockAccountsService).GetAccountByIDFunc = func(ctx context.Context, id uint) (*entities.Account, error) {
					return testAccount, nil
				}
				controller.AccountsService.(*MockAccountsService).GetAccountTransactionsFunc = func(ctx context.Context, in repository.GetTransactionsIn) ([]*entities.Transaction, error) {
					// Проверяем что время парсится правильно
					assert.Equal(t, time.UTC, in.DateFrom.Location())
					assert.Equal(t, 10, in.DateFrom.Hour())

					_, offset := in.DateTo.Zone()
					assert.Equal(t, 5*3600, offset) // +05:00
					assert.Equal(t, 23, in.DateTo.Hour())

					return testTransactions[:1], nil
				}
			},
			expectedStatus: 200,
			expectedCount:  1,
		},
		{
			name: "аккаунт не найден",
			url:  server.URL + "/accounts/999/transactions",
			mockSetup: func() {
				controller.AccountsService.(*MockAccountsService).GetAccountByIDFunc = func(ctx context.Context, id uint) (*entities.Account, error) {
					return nil, fmt.Errorf("account not found")
				}
			},
			expectedStatus: 404,
			shouldHaveErr:  true,
		},
		{
			name: "некорректный формат даты",
			url:  server.URL + "/accounts/123/transactions?date_from=invalid-date",
			mockSetup: func() {
				controller.AccountsService.(*MockAccountsService).GetAccountByIDFunc = func(ctx context.Context, id uint) (*entities.Account, error) {
					return testAccount, nil
				}
			},
			expectedStatus: 400,
			shouldHaveErr:  true,
		},
		{
			name: "доступ к чужому аккаунту",
			url:  server.URL + "/accounts/123/transactions",
			mockSetup: func() {
				controller.AccountsService.(*MockAccountsService).GetAccountByIDFunc = func(ctx context.Context, id uint) (*entities.Account, error) {
					// Возвращаем аккаунт с другим customer ID
					return &entities.Account{
						ID:         123,
						CustomerID: 999, // другой customer
						Balance:    types.NewMoney(decimal.NewFromFloat(1000.0), types.Currency("UZS")),
					}, nil
				}
			},
			expectedStatus: 403,
			shouldHaveErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Настройка моков
			tt.mockSetup()

			// Выполняем HTTP запрос
			resp, err := http.Get(tt.url)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Проверяем статус
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Декодируем ответ
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			if tt.shouldHaveErr {
				assert.Contains(t, response, "error")
			} else {
				assert.Contains(t, response, "transactions")
				transactions, ok := response["transactions"].([]interface{})
				assert.True(t, ok)
				assert.Len(t, transactions, tt.expectedCount)
			}
		})
	}
}

func TestController_GetAccountTransactions_HTTPServer_DateParsing(t *testing.T) {
	server, controller := setupHTTPTestServer()
	defer server.Close()

	testAccount := &entities.Account{
		ID:         123,
		CustomerID: 456,
		Balance:    types.NewMoney(decimal.NewFromFloat(1000.0), types.Currency("UZS")),
	}

	tests := []struct {
		name        string
		queryParams string
		expectError bool
		desc        string
	}{
		{
			name:        "простые даты в RFC3339",
			queryParams: "date_from=2024-01-15T00:00:00Z&date_to=2024-01-25T23:59:59Z",
			expectError: false,
			desc:        "должен парсить RFC3339 формат",
		},
		{
			name:        "время с UTC",
			queryParams: "date_from=2024-01-15T10:30:00Z&date_to=2024-01-25T15:45:00Z",
			expectError: false,
			desc:        "должен парсить время в UTC",
		},
		{
			name:        "время с timezone offset",
			queryParams: "date_from=2024-01-15T10:30:00%2B05:00&date_to=2024-01-25T15:45:00-03:00",
			expectError: false,
			desc:        "должен парсить время с различными timezone",
		},
		{
			name:        "некорректный формат",
			queryParams: "date_from=not-a-date&date_to=2024-01-25T00:00:00Z",
			expectError: true,
			desc:        "должен возвращать ошибку для некорректного формата",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedParams repository.GetTransactionsIn

			controller.AccountsService.(*MockAccountsService).GetAccountByIDFunc = func(ctx context.Context, id uint) (*entities.Account, error) {
				return testAccount, nil
			}

			controller.AccountsService.(*MockAccountsService).GetAccountTransactionsFunc = func(ctx context.Context, in repository.GetTransactionsIn) ([]*entities.Transaction, error) {
				capturedParams = in
				return []*entities.Transaction{}, nil
			}

			url := server.URL + "/accounts/123/transactions?" + tt.queryParams
			resp, err := http.Get(url)
			require.NoError(t, err)
			defer resp.Body.Close()

			if tt.expectError {
				assert.Equal(t, 400, resp.StatusCode, tt.desc)
			} else {
				assert.Equal(t, 200, resp.StatusCode, tt.desc)

				// Проверяем что параметры правильно переданы
				assert.Equal(t, uint(123), capturedParams.AccountID)

				// Если были переданы даты, проверяем что они не пустые
				if !capturedParams.DateFrom.IsZero() {
					assert.False(t, capturedParams.DateFrom.IsZero(), "DateFrom должен быть установлен")
				}
				if !capturedParams.DateTo.IsZero() {
					assert.False(t, capturedParams.DateTo.IsZero(), "DateTo должен быть установлен")
				}
			}
		})
	}
}
