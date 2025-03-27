package mongodb

import (
	"context"
	"time"

	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/repository"
	"github.com/bekha-io/openbank/domain/types"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoTransaction struct {
	ID            uint                 `bson:"id"`
	ExternalID    string               `bson:"external_id"`
	FromAccountId uint                 `bson:"from_account_id"`
	ToAccountId   uint                 `bson:"to_account_id"`
	Category      string               `bson:"category"`
	Status        string               `bson:"status"`
	StatusReason  string               `bson:"status_reason"`
	Comment       string               `bson:"comment"`
	Amount        primitive.Decimal128 `bson:"amount"`
	Currency      string               `bson:"currency"`
	CreatedAt     time.Time            `bson:"created_at"`
}

func (m *mongoTransaction) ParseEntity(e *entities.Transaction) {
	m.ID = e.ID
	m.FromAccountId = e.FromAccountId
	m.ExternalID = e.ExternalID
	m.ToAccountId = e.ToAccountId
	m.Category = string(e.Category)
	m.Status = string(e.Status)
	m.StatusReason = e.StatusReason
	m.Comment = e.Comment
	m.Amount, _ = primitive.ParseDecimal128(e.Amount.StringFixed(2))
	m.CreatedAt = e.CreatedAt
	// если нужна поддержка валют — добавь
}

func (m *mongoTransaction) ToEntity() *entities.Transaction {
	amount, _ := decimal.NewFromString(m.Amount.String())

	return &entities.Transaction{
		ID:            m.ID,
		ExternalID:    m.ExternalID,
		FromAccountId: m.FromAccountId,
		ToAccountId:   m.ToAccountId,
		Category:      types.TransactionCategory(m.Category),
		Status:        entities.TransactionStatus(m.Status),
		StatusReason:  m.StatusReason,
		Comment:       m.Comment,
		Amount:        amount,
		CreatedAt:     m.CreatedAt,
	}
}

var _ repository.ITransactionRepository = (*MongoTransactionRepository)(nil)

type MongoTransactionRepository struct {
	dbName string
	cl     *mongo.Client
}

// GetTransactionsByAccountID implements repository.ITransactionRepository.
func (r *MongoTransactionRepository) GetTransactionsByAccountID(ctx context.Context, id uint) ([]*entities.Transaction, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"from_account_id": id},
			{"to_account_id": id},
		},
	}

	cur, err := r.cl.Database(r.dbName).Collection("transactions").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var mongoTransactions []*mongoTransaction
	if err := cur.All(ctx, &mongoTransactions); err != nil {
		return nil, err
	}

	var transactions []*entities.Transaction
	for _, mt := range mongoTransactions {
		transactions = append(transactions, mt.ToEntity())
	}
	return transactions, nil
}

// SaveTransaction implements repository.ITransactionRepository.
func (r *MongoTransactionRepository) SaveTransaction(ctx context.Context, tr *entities.Transaction) error {
	doc := &mongoTransaction{}
	doc.ParseEntity(tr)

	_, err := r.cl.Database(r.dbName).Collection("transactions").UpdateOne(
		ctx,
		bson.M{"id": doc.ID},
		bson.D{{Key: "$set", Value: doc}},
		options.Update().SetUpsert(true),
	)
	return err
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
func (r *MongoTransactionRepository) GetByID(ctx context.Context, id uint) (*entities.Transaction, error) {
	return r.GetBy(ctx, "id", id)
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
