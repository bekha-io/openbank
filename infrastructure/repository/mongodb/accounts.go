package mongodb

import (
	"context"
	"log/slog"

	"github.com/bekha-io/vaultonomy/domain/entities"
	"github.com/bekha-io/vaultonomy/domain/repository"
	"github.com/bekha-io/vaultonomy/domain/types"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoAccount struct {
	ID         string               `bson:"id"`
	CustomerID string               `bson:"customer_id"`
	Balance    primitive.Decimal128 `bson:"balance"`
	Currency   string               `bson:"currency"`
}

func (c *mongoAccount) ParseEntity(e *entities.Account) {
	c.ID = string(e.ID)
	c.Balance, _ = primitive.ParseDecimal128(e.Balance.Amount.StringFixed(2))
	c.Currency = string(e.Balance.Currency)
	c.CustomerID = e.CustomerID.String()
}

func (c *mongoAccount) ToEntity() *entities.Account {
	return &entities.Account{
		ID:         types.AccountID(c.ID),
		CustomerID: types.CustomerID(uuid.MustParse(c.CustomerID)),
		Balance:    types.NewMoney(decimal.RequireFromString(c.Balance.String()), types.Currency(c.Currency)),
	}
}

var _ repository.IAccountRepository = (*MongoAccountRepository)(nil)

type MongoAccountRepository struct {
	cl     *mongo.Client
	dbName string
}

func NewMongoAccountRepository(cl *mongo.Client, dbName string) repository.IAccountRepository {
	return &MongoAccountRepository{
		cl:     cl,
		dbName: dbName,
	}
}

// GetBy implements repository.IAccountRepository.
func (r *MongoAccountRepository) GetBy(ctx context.Context, key string, value interface{}) (*entities.Account, error) {
	var row mongoAccount
	err := r.cl.Database(r.dbName).Collection("accounts").
		FindOne(ctx, bson.M{key: value}).Decode(&row)
	if err != nil {
		return nil, err
	}
	return row.ToEntity(), nil
}

// GetByID implements repository.IAccountRepository.
func (r *MongoAccountRepository) GetByID(ctx context.Context, id types.AccountID) (*entities.Account, error) {
	return r.GetBy(ctx, "id", string(id))
}

// Save implements repository.IAccountRepository.
func (r *MongoAccountRepository) Save(ctx context.Context, acc *entities.Account) error {
	a := &mongoAccount{}
	a.ParseEntity(acc)

	_, err := r.cl.Database(r.dbName).Collection("accounts").UpdateOne(ctx, bson.M{"id": a.ID},
		bson.D{{Key: "$set", Value: a}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

// GetManyBy implements repository.IAccountRepository.
func (r *MongoAccountRepository) GetManyBy(ctx context.Context, filters ...repository.Filter) ([]*entities.Account, error) {
	var filter = bson.D{}
	for _, flt := range filters {
		filter = append(filter, ParseFilter(flt))
	}
	slog.Info("query", "map", filter)

	cur, err := r.cl.Database(r.dbName).Collection("accounts").Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var mongoAccounts []*mongoAccount
	err = cur.All(ctx, &mongoAccounts)
	if err != nil {
		return nil, err
	}

	var accounts []*entities.Account
	for _, mongoAcc := range mongoAccounts {
		accounts = append(accounts, mongoAcc.ToEntity())
	}

	return accounts, nil
}
