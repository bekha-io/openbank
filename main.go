package main

import (
	"context"
	"os"

	"github.com/bekha-io/vaultonomy/domain/services"
	"github.com/bekha-io/vaultonomy/infrastructure/repository/mongodb"
	"github.com/bekha-io/vaultonomy/presentation/rest/accounts"
	"github.com/bekha-io/vaultonomy/presentation/rest/customers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	mongoUri := os.Getenv("MONGODB_URI")
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}

	transactionsRepo := mongodb.NewMongoTransactionRepository(mongoClient, "vaultonomy")
	accountsRepo := mongodb.NewMongoAccountRepository(mongoClient, "vaultonomy")
	individualCustomersRepo := mongodb.NewMongoIndividualCustomerRepository(mongoClient, "vaultonomy")

	accountsSvc := services.NewAccountsService(accountsRepo, transactionsRepo)
	individualCustomersSvc := services.NewIndividualCustomerService(individualCustomersRepo, accountsRepo)

	accountsController := accounts.NewAccountsController(accountsSvc)
	customersController := customers.NewCustomerController(individualCustomersSvc)

	r := gin.Default()

	// api/v1
	v1 := r.Group("/v1")
	{
		// api/v1/customers
		customersGroup := v1.Group("/customers")
		{
			customersGroup.POST("", customersController.CreateCustomer)
			customersGroup.GET("/:id/accounts", customersController.GetCustomerAccounts)
		}

		accountsGroup := v1.Group("/accounts")
		{
			accountsGroup.POST("", accountsController.CreateAccount)
			accountsGroup.POST("/:id/withdraw", accountsController.Withdraw)
			accountsGroup.POST("/:id/deposit", accountsController.Deposit)
		}
	}

	r.Run(":"+os.Getenv("APP_PORT"))
}
