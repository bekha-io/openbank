package mongodb

import (
	"context"

	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/repository"
	"github.com/bekha-io/openbank/domain/types"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoIndividualCustomer struct {
	ID          string `bson:"id"`
	PhoneNumber string `bson:"phone_number"`
	FirstName   string `bson:"first_name"`
	LastName    string `bson:"last_name"`
	MiddleName  string `bson:"middle_name"`
}

func (c *mongoIndividualCustomer) ParseEntity(e *entities.IndividualCustomer) {
	c.ID = e.ID.String()
	c.PhoneNumber = e.PhoneNumber
	c.FirstName = e.FirstName
	c.LastName = e.LastName
	c.MiddleName = e.MiddleName
}

func (c *mongoIndividualCustomer) ToEntity() *entities.IndividualCustomer {
	return &entities.IndividualCustomer{
		ID:          types.CustomerID(uuid.MustParse(c.ID)),
		PhoneNumber: c.PhoneNumber,
		FirstName:   c.FirstName,
		MiddleName:  c.MiddleName,
		LastName:    c.LastName,
	}
}

var _ repository.IIndividualCustomerRepository = (*MongoIndividualCustomerRepository)(nil)

type MongoIndividualCustomerRepository struct {
	cl     *mongo.Client
	dbName string
}

func NewMongoIndividualCustomerRepository(cl *mongo.Client, dbName string) repository.IIndividualCustomerRepository {
	return &MongoIndividualCustomerRepository{
		cl:     cl,
		dbName: dbName,
	}
}

// GetBy implements repository.IIndividualCustomerRepository.
func (r *MongoIndividualCustomerRepository) GetBy(ctx context.Context, key string, value interface{}) (*entities.IndividualCustomer, error) {
	var row mongoIndividualCustomer
	err := r.cl.Database(r.dbName).Collection("individual_customers").
		FindOne(ctx, bson.M{key: value}).Decode(&row)
	if err != nil {
		return nil, err
	}
	return row.ToEntity(), nil
}

// GetByID implements repository.IIndividualCustomerRepository.
func (r *MongoIndividualCustomerRepository) GetByID(ctx context.Context, id types.CustomerID) (*entities.IndividualCustomer, error) {
	return r.GetBy(ctx, "id", id.String())
}

// Save implements repository.IIndividualCustomerRepository.
func (r *MongoIndividualCustomerRepository) Save(ctx context.Context, customer *entities.IndividualCustomer) error {
	c := &mongoIndividualCustomer{}
	c.ParseEntity(customer)

	_, err := r.cl.Database(r.dbName).Collection("individual_customers").UpdateOne(ctx, bson.M{"id": customer.ID.String()},
		bson.D{{Key: "$set", Value: c}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}
