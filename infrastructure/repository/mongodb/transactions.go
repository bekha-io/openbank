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

type mongoTransaction struct {
	ID              string               `bson:"id"`
	AccountID       string               `bson:"account_id"`
	TransactionType string               `bson:"transaction_type"`
	Amount          primitive.Decimal128 `bson:"amount"`
	Currency        string               `bson:"currency"`
	CreatedAt       time.Time            `bson:"created_at"`
	Comment         string               `bson:"comment"`
}

func (c *mongoTransaction) ParseEntity(e *entities.Transaction) {
	c.ID = e.ID.String()
	c.Amount, _ = primitive.ParseDecimal128(e.Amount.Amount.StringFixed(2))
	c.Currency = string(e.Amount.Currency)
	c.AccountID = string(e.AccountID)
	c.CreatedAt = e.CreatedAt
	c.TransactionType = string(e.TransactionType)
	c.Comment = e.Comment
}

func (c *mongoTransaction) ToEntity() *entities.Transaction {
	return &entities.Transaction{
		ID:              types.TransactionID(uuid.MustParse(c.ID)),
		AccountID:       types.AccountID(c.AccountID),
		TransactionType: types.TransactionType(c.TransactionType),
		Amount:          types.NewMoney(decimal.RequireFromString(c.Amount.String()), types.Currency(c.Currency)),
		Comment:         c.Comment,
	}
}

var _ repository.ITransactionRepository = (*MongoTransactionRepository)(nil)

type MongoTransactionRepository struct {
	dbName string
	cl     *mongo.Client
}

func NewMongoTransactionRepository(cl *mongo.Client, dbName string) *MongoTransactionRepository {
	return &MongoTransactionRepository{
		dbName: dbName,
		cl:     cl,
	}
}

// GetBy implements repository.ITransactionRepository.
func (r *MongoTransactionRepository) GetBy(ctx context.Context, key string, value interface{}) (*entities.Transaction, error) {
	var row mongoTransaction
	err := r.cl.Database(r.dbName).Collection("transactions").
		FindOne(ctx, bson.M{key: value}).Decode(&row)
	if err != nil {
		return nil, err
	}
	return row.ToEntity(), nil
}

// GetByID implements repository.ITransactionRepository.
func (r *MongoTransactionRepository) GetByID(ctx context.Context, id types.TransactionID) (*entities.Transaction, error) {
	return r.GetBy(ctx, "id", id.String())
}

// GetManyBy implements repository.ITransactionRepository.
func (r *MongoTransactionRepository) GetManyBy(ctx context.Context, filters ...repository.Filter) ([]*entities.Transaction, error) {
	var filter = bson.D{}
	for _, flt := range filters {
		filter = append(filter, ParseFilter(flt))
	}

	cur, err := r.cl.Database(r.dbName).Collection("transactions").Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var mongoTransactions []*mongoTransaction
	err = cur.All(ctx, &mongoTransactions)
	if err != nil {
		return nil, err
	}

	var transactions []*entities.Transaction
	for _, mongoTransaction := range mongoTransactions {
		transactions = append(transactions, mongoTransaction.ToEntity())
	}

	return transactions, nil
}

// Save implements repository.ITransactionRepository.
func (r *MongoTransactionRepository) Save(ctx context.Context, transaction *entities.Transaction) error {
	a := &mongoTransaction{}
	a.ParseEntity(transaction)

	_, err := r.cl.Database(r.dbName).Collection("transactions").UpdateOne(ctx, bson.M{"id": a.ID},
		bson.D{{Key: "$set", Value: a}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}
