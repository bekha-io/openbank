package mongodb

import (
	"context"

	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/repository"
	"github.com/bekha-io/openbank/domain/types"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (c *mongoIndividualCustomer) ParseEntity(e *entities.Customer) {
	c.ID = e.ID.String()
	c.PhoneNumber = e.PhoneNumber
	c.FirstName = e.FirstName
	c.LastName = e.LastName
	c.MiddleName = e.MiddleName
}

func (c *mongoIndividualCustomer) ToEntity() *entities.Customer {
	return &entities.Customer{
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
func (r *MongoIndividualCustomerRepository) GetBy(ctx context.Context, key string, value interface{}) (*entities.Customer, error) {
	var row mongoIndividualCustomer
	err := r.cl.Database(r.dbName).Collection("individual_customers").
		FindOne(ctx, bson.M{key: value}).Decode(&row)
	if err != nil {
		return nil, err
	}
	return row.ToEntity(), nil
}

// GetByID implements repository.IIndividualCustomerRepository.
func (r *MongoIndividualCustomerRepository) GetByID(ctx context.Context, id types.CustomerID) (*entities.Customer, error) {
	return r.GetBy(ctx, "id", id.String())
}

// Save implements repository.IIndividualCustomerRepository.
func (r *MongoIndividualCustomerRepository) Save(ctx context.Context, customer *entities.Customer) error {
	c := &mongoIndividualCustomer{}
	c.ParseEntity(customer)

	_, err := r.cl.Database(r.dbName).Collection("individual_customers").UpdateOne(ctx, bson.M{"id": customer.ID.String()},
		bson.D{{Key: "$set", Value: c}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

// GetManyIDLike implements repository.IIndividualCustomerRepository.
func (r *MongoIndividualCustomerRepository) GetManyIDLike(ctx context.Context, id types.Currency) ([]*entities.Customer, error) {
	pattern := ".*" + id + ".*"
	cur, err := r.cl.Database(r.dbName).Collection("individual_customers").
		Find(ctx, bson.M{"id": bson.M{"$regex": primitive.Regex{Pattern: string(pattern)}}})
	if err != nil {
		return nil, err
	}

	var mongoCustomers []*mongoIndividualCustomer
	err = cur.All(ctx, &mongoCustomers)
	if err != nil {
		return nil, err
	}

	var customers []*entities.Customer
	for _, mongoAcc := range mongoCustomers {
		customers = append(customers, mongoAcc.ToEntity())
	}

	return customers, nil
}

func (r *MongoIndividualCustomerRepository) GetManyPhoneNumberLike(ctx context.Context, phoneNumber string) ([]*entities.Customer, error) {
	pattern := ".*" + phoneNumber + ".*"
	cur, err := r.cl.Database(r.dbName).Collection("individual_customers").
		Find(ctx, bson.M{"phone_number": bson.M{"$regex": primitive.Regex{Pattern: string(pattern)}}})
	if err != nil {
		return nil, err
	}

	var mongoCustomers []*mongoIndividualCustomer
	err = cur.All(ctx, &mongoCustomers)
	if err != nil {
		return nil, err
	}

	var customers []*entities.Customer
	for _, mongoAcc := range mongoCustomers {
		customers = append(customers, mongoAcc.ToEntity())
	}

	return customers, nil
}
