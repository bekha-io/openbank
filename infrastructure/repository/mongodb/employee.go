package mongodb

import (
	"context"
	"time"

	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/domain/repository"
	"github.com/bekha-io/openbank/domain/types"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoEmployee struct {
	ID         string    `bson:"id"`
	Email      string    `bson:"email"`
	Password   string    `bson:"password"`
	FirstName  string    `bson:"first_name"`
	LastName   string    `bson:"last_name"`
	MiddleName string    `bson:"middle_name"`
	CreatedAt  time.Time `bson:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at"`
}

func (e *mongoEmployee) ToEntity() *entities.Employee {
	return &entities.Employee{
		ID:         types.EmployeeID(uuid.MustParse(e.ID)),
		Email:      e.Email,
		Password:   e.Password,
		FirstName:  e.FirstName,
		LastName:   e.LastName,
		MiddleName: e.MiddleName,
		CreatedAt:  e.CreatedAt,
	}
}

func (em *mongoEmployee) ParseEntity(e *entities.Employee) {
	em.ID = e.ID.String()
	em.Email = e.Email
	em.Password = e.Password
	em.FirstName = e.FirstName
	em.LastName = e.LastName
	em.MiddleName = e.MiddleName
	em.CreatedAt = e.CreatedAt
	em.UpdatedAt = e.UpdatedAt
}

var _ repository.IEmployeeRepository = (*MongoEmployeeRepository)(nil)

type MongoEmployeeRepository struct {
	cl     *mongo.Client
	dbName string
}

func NewMongoEmployeeRepository(client *mongo.Client, dbName string) *MongoEmployeeRepository {
	return &MongoEmployeeRepository{client, dbName}
}

func (r *MongoEmployeeRepository) GetByEmail(ctx context.Context, email string) (*entities.Employee, error) {
	var result mongoEmployee
	err := r.cl.Database(r.dbName).Collection("employees").FindOne(ctx, bson.M{"email": email}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return result.ToEntity(), err
}

func (r *MongoEmployeeRepository) GetManyLike(ctx context.Context, query string) ([]*entities.Employee, error) {
	pattern := ".*" + query + ".*"
	regexFilter := bson.M{"$regex": primitive.Regex{Pattern: string(pattern)}}
	filter := bson.M{"$or": []bson.M{bson.M{"email": regexFilter}, bson.M{"first_name": regexFilter}, bson.M{"last_name": regexFilter}, bson.M{"middle_name": regexFilter}}}

	var results []*entities.Employee
	cursor, err := r.cl.Database(r.dbName).Collection("employees").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result *mongoEmployee
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result.ToEntity())
	}

	return results, nil
}

func (r *MongoEmployeeRepository) Save(ctx context.Context, employee *entities.Employee) error {
	var mongoEmp mongoEmployee
	mongoEmp.ParseEntity(employee)
	_, err := r.cl.Database(r.dbName).Collection("employees").UpdateOne(ctx, bson.M{"id": mongoEmp.ID},
		bson.D{{Key: "$set", Value: mongoEmp}}, options.Update().SetUpsert(true))
	return err
}
