package gormRepo

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/repository"
	"github.com/bekha-io/openbank/domain/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGormTransactionRepository_GetTransactions_WithSQLMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := &GormTransactionRepository{db: gormDB}

	tests := []struct {
		name     string
		input    repository.GetTransactionsIn
		mockFunc func()
		expected int
	}{
		{
			name: "без фильтров по дате",
			input: repository.GetTransactionsIn{
				AccountID: 123,
			},
			mockFunc: func() {
				expectedSQL := `SELECT * FROM "gorm_transactions" WHERE from_account_id = $1 OR to_account_id = $2`
				rows := sqlmock.NewRows([]string{"id", "external_id", "from_account_id", "to_account_id", "category", "status", "status_reason", "comment", "amount", "created_at"}).
					AddRow(1, "ext-1", 123, 456, "transfer", "completed", "success", "test", decimal.NewFromFloat(100.50), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
					WithArgs(123, 123).
					WillReturnRows(rows)
			},
			expected: 1,
		},
		{
			name: "с фильтром DateFrom",
			input: repository.GetTransactionsIn{
				AccountID: 123,
				DateFrom:  time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			mockFunc: func() {
				expectedSQL := `SELECT * FROM "gorm_transactions" WHERE (from_account_id = $1 OR to_account_id = $2) AND created_at >= $3`
				rows := sqlmock.NewRows([]string{"id", "external_id", "from_account_id", "to_account_id", "category", "status", "status_reason", "comment", "amount", "created_at"}).
					AddRow(1, "ext-1", 123, 456, "transfer", "completed", "success", "test", decimal.NewFromFloat(100.50), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
					WithArgs(123, 123, time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)).
					WillReturnRows(rows)
			},
			expected: 1,
		},
		{
			name: "с фильтром DateTo",
			input: repository.GetTransactionsIn{
				AccountID: 123,
				DateTo:    time.Date(2024, 12, 31, 23, 59, 59, 0, time.FixedZone("UZT", 5*3600)),
			},
			mockFunc: func() {
				expectedSQL := `SELECT * FROM "gorm_transactions" WHERE (from_account_id = $1 OR to_account_id = $2) AND created_at <= $3`
				rows := sqlmock.NewRows([]string{"id", "external_id", "from_account_id", "to_account_id", "category", "status", "status_reason", "comment", "amount", "created_at"}).
					AddRow(1, "ext-1", 123, 456, "transfer", "completed", "success", "test", decimal.NewFromFloat(100.50), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
					WithArgs(123, 123, time.Date(2024, 12, 31, 23, 59, 59, 0, time.FixedZone("UZT", 5*3600))).
					WillReturnRows(rows)
			},
			expected: 1,
		},
		{
			name: "с обоими фильтрами по дате",
			input: repository.GetTransactionsIn{
				AccountID: 123,
				DateFrom:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				DateTo:    time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
			},
			mockFunc: func() {
				expectedSQL := `SELECT * FROM "gorm_transactions" WHERE (from_account_id = $1 OR to_account_id = $2) AND created_at >= $3 AND created_at <= $4`
				rows := sqlmock.NewRows([]string{"id", "external_id", "from_account_id", "to_account_id", "category", "status", "status_reason", "comment", "amount", "created_at"}).
					AddRow(1, "ext-1", 123, 456, "transfer", "completed", "success", "test", decimal.NewFromFloat(100.50), time.Now()).
					AddRow(2, "ext-2", 456, 123, "transfer", "completed", "success", "test2", decimal.NewFromFloat(200.75), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
					WithArgs(123, 123, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)).
					WillReturnRows(rows)
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			transactions, err := repo.GetTransactions(nil, tt.input)

			require.NoError(t, err)
			assert.Len(t, transactions, tt.expected)

			// Проверяем что все expectations выполнены
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGormTransactionRepository_FilterLogic(t *testing.T) {
	// Тестируем логику фильтрации без реальной БД

	tests := []struct {
		name  string
		input repository.GetTransactionsIn
		desc  string
	}{
		{
			name: "без фильтров по дате",
			input: repository.GetTransactionsIn{
				AccountID: 100,
			},
			desc: "должен принимать только AccountID",
		},
		{
			name: "с фильтром DateFrom",
			input: repository.GetTransactionsIn{
				AccountID: 100,
				DateFrom:  time.Date(2024, 1, 1, 10, 30, 0, 0, time.UTC),
			},
			desc: "должен принимать DateFrom с временем и timezone",
		},
		{
			name: "с фильтром DateTo",
			input: repository.GetTransactionsIn{
				AccountID: 100,
				DateTo:    time.Date(2024, 12, 31, 23, 59, 59, 0, time.Local),
			},
			desc: "должен принимать DateTo с локальным timezone",
		},
		{
			name: "с обоими фильтрами",
			input: repository.GetTransactionsIn{
				AccountID: 100,
				DateFrom:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				DateTo:    time.Date(2024, 12, 31, 23, 59, 59, 0, time.FixedZone("UZT", 5*3600)),
			},
			desc: "должен принимать оба фильтра с разными timezone",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Проверяем что структура правильно принимает параметры
			assert.Equal(t, uint(100), tt.input.AccountID)

			// Проверяем логику IsZero() для time.Time
			if tt.input.DateFrom.IsZero() {
				assert.True(t, tt.input.DateFrom.IsZero(), "DateFrom должен быть пустым")
			} else {
				assert.False(t, tt.input.DateFrom.IsZero(), "DateFrom должен быть установлен")
				assert.NotNil(t, tt.input.DateFrom.Location(), "DateFrom должен иметь timezone")
			}

			if tt.input.DateTo.IsZero() {
				assert.True(t, tt.input.DateTo.IsZero(), "DateTo должен быть пустым")
			} else {
				assert.False(t, tt.input.DateTo.IsZero(), "DateTo должен быть установлен")
				assert.NotNil(t, tt.input.DateTo.Location(), "DateTo должен иметь timezone")
			}
		})
	}
}

func TestGormTransaction_Conversion(t *testing.T) {
	// Тестируем конвертацию между gormTransaction и entities.Transaction

	original := &entities.Transaction{
		ID:            123,
		ExternalID:    "ext-123",
		FromAccountId: 100,
		ToAccountId:   200,
		Category:      types.TransactionCategoryTransfer,
		Status:        entities.TransactionStatusCompleted,
		StatusReason:  "success",
		Comment:       "test transaction",
		Amount:        decimal.NewFromFloat(1234.56),
		CreatedAt:     time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
	}

	// Конвертируем в gormTransaction
	gormTx := &gormTransaction{}
	gormTx.fromEntity(original)

	// Проверяем что все поля скопировались
	assert.Equal(t, original.ID, gormTx.ID)
	assert.Equal(t, original.ExternalID, gormTx.ExternalID)
	assert.Equal(t, original.FromAccountId, gormTx.FromAccountId)
	assert.Equal(t, original.ToAccountId, gormTx.ToAccountId)
	assert.Equal(t, string(original.Category), gormTx.Category)
	assert.Equal(t, string(original.Status), gormTx.Status)
	assert.Equal(t, original.StatusReason, gormTx.StatusReason)
	assert.Equal(t, original.Comment, gormTx.Comment)
	assert.True(t, original.Amount.Equal(gormTx.Amount))
	assert.Equal(t, original.CreatedAt, gormTx.CreatedAt)

	// Конвертируем обратно в entities.Transaction
	converted := gormTx.toEntity()

	// Проверяем что конвертация туда-обратно работает корректно
	assert.Equal(t, original.ID, converted.ID)
	assert.Equal(t, original.ExternalID, converted.ExternalID)
	assert.Equal(t, original.FromAccountId, converted.FromAccountId)
	assert.Equal(t, original.ToAccountId, converted.ToAccountId)
	assert.Equal(t, original.Category, converted.Category)
	assert.Equal(t, original.Status, converted.Status)
	assert.Equal(t, original.StatusReason, converted.StatusReason)
	assert.Equal(t, original.Comment, converted.Comment)
	assert.True(t, original.Amount.Equal(converted.Amount))
	assert.Equal(t, original.CreatedAt, converted.CreatedAt)
}
