package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/bekha-io/openbank/domain/entities"
	"github.com/bekha-io/openbank/infrastructure/repository/mongodb"
	"github.com/sethvargo/go-password/password"
	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func CreateEmployee(c *cli.Context) error {
	email := c.String("email")
	if email == "" {
		return errors.New("email is required")
	}

	pwd := c.String("password")
	if pwd == "" {
		pwd, _ = password.Generate(6, 0, 0, false, true)
	}

	firstName := c.String("first-name")
	lastName := c.String("last-name")
	middleName := c.String("middle-name")

	employee := entities.NewEmployee(email, pwd, firstName, lastName, middleName)

	employeesRepo := mongodb.NewMongoEmployeeRepository(MongoClient, "openbank")
	err := employeesRepo.Save(c.Context, employee)
	if err!= nil {
        return err
    }

	log.Printf("Created! Email: %v | Password: %v \n", employee.Email, employee.Password)
	return nil
}


func main() {
	var err error
	mongoUri := os.Getenv("MONGODB_URI")
	MongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}

	app := &cli.App{
		Name: "openbank-cli",
		Commands: []*cli.Command{
			{
				Name: "customers",
				Subcommands: []*cli.Command{
					{
						Name: "create",
						Action: CreateEmployee,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}