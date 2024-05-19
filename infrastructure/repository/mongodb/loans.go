package mongodb

import (
	"context"
	"time"

	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/repository"
	"github.com/bekha-io/openbank/domain/types"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ repository.ILoanRepository = (*MongoLoanRepository)(nil)

type mongoLoanProduct struct {
	ID                  string               `bson:"id"`
	Name                string               `bson:"name"`
	MinDuration         uint                  `bson:"min_duration"`
	MaxDuration         uint                  `bson:"max_duration"`
	Currency            string               `bson:"currency"`
	MinAmount           primitive.Decimal128 `bson:"min_amount"`
	MaxAmount           primitive.Decimal128 `bson:"max_amount"`
	InterestRate        primitive.Decimal128 `bson:"interest_rate"`
	LoanType            string               `bson:"loan_type"`
	DailyOverduePenalty primitive.Decimal128 `bson:"daily_overdue_penalty"`
	CreatedAt           time.Time            `bson:"created_at"`
	UpdatedAt           time.Time            `bson:"updated_at"`
}

func (m *mongoLoanProduct) ParseEntity(e *entities.LoanProduct) {
	m.ID = e.ID.String()
	m.Name = e.Name
	m.MinDuration = e.MinDuration
	m.MaxDuration = e.MaxDuration
	m.Currency = string(e.MinAmount.Currency)
	m.MaxAmount, _ = primitive.ParseDecimal128(e.MaxAmount.Amount.StringFixed(2))
	m.MinAmount, _ = primitive.ParseDecimal128(e.MinAmount.Amount.StringFixed(2))
	m.InterestRate, _ = primitive.ParseDecimal128(e.InterestRate.StringFixed(2))
	m.LoanType = string(e.LoanType)
	m.DailyOverduePenalty, _ = primitive.ParseDecimal128(e.DailyOverduePenalty.StringFixed(2))
	m.CreatedAt = e.CreatedAt
	m.UpdatedAt = e.UpdatedAt
}

func (m *mongoLoanProduct) ToEntity() *entities.LoanProduct {
	var e entities.LoanProduct

	e.ID = types.LoanProductID(uuid.MustParse(m.ID))
	e.Name = m.Name
	e.MinDuration = m.MinDuration
	e.MaxDuration = m.MaxDuration
	e.MinAmount = *types.NewMoney(decimal.RequireFromString(m.MinAmount.String()), types.Currency(m.Currency))
	e.MaxAmount = *types.NewMoney(decimal.RequireFromString(m.MaxAmount.String()), types.Currency(m.Currency))
	e.InterestRate = decimal.RequireFromString(m.InterestRate.String())
	e.LoanType = types.LoanType(m.LoanType)
	e.DailyOverduePenalty = decimal.RequireFromString(m.DailyOverduePenalty.String())
	e.CreatedAt = m.CreatedAt
	e.UpdatedAt = m.UpdatedAt

	return &e
}

var _ repository.ILoanRepository = (*MongoLoanRepository)(nil)

type MongoLoanRepository struct {
	cl     *mongo.Client
	dbName string
}

func NewMongoLoanRepository(cl *mongo.Client, dbName string) repository.ILoanRepository {
	return &MongoLoanRepository{
		cl:     cl,
		dbName: dbName,
	}
}

// GetLoanProductByID implements repository.ILoanRepository.
func (m *MongoLoanRepository) GetLoanProductByID(ctx context.Context, id types.LoanProductID) (*entities.LoanProduct, error) {
	var result *mongoLoanProduct
	err := m.cl.Database(m.dbName).Collection("loan_products").FindOne(ctx, bson.M{"id": id.String()}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.ToEntity(), nil
}

// SaveLoanProduct implements repository.ILoanRepository.
func (m *MongoLoanRepository) SaveLoanProduct(ctx context.Context, lp *entities.LoanProduct) error {
	var result = &mongoLoanProduct{}
	result.ParseEntity(lp)
	_, err := m.cl.Database(m.dbName).Collection("loan_products").UpdateOne(ctx, bson.M{"id": result.ID}, bson.D{
		{Key: "$set", Value: result},
	},
		options.Update().SetUpsert(true))
	return err
}

// GetAllLoanProducts implements repository.ILoanRepository.
func (m *MongoLoanRepository) GetAllLoanProducts(ctx context.Context) ([]*entities.LoanProduct, error) {
	var results []*entities.LoanProduct
	cursor, err := m.cl.Database(m.dbName).Collection("loan_products").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result *mongoLoanProduct
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result.ToEntity())
	}

	return results, nil
}
