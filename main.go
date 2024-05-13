package main

import (
	"context"
	"os"

	"github.com/bekha-io/openbank/domain/services"
	"github.com/bekha-io/openbank/infrastructure/repository/mongodb"
	"github.com/bekha-io/openbank/presentation/rest/accounts"
	"github.com/bekha-io/openbank/presentation/rest/customers"
	"github.com/bekha-io/openbank/presentation/rest/loans"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	decimal.MarshalJSONWithoutQuotes = true
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

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
	loansSvc := services.NewLoanService()

	accountsController := accounts.NewAccountsController(accountsSvc)
	customersController := customers.NewCustomerController(individualCustomersSvc)
	loansController := loans.NewLoanController(loansSvc)

	r := gin.Default()
	r.Use(CORSMiddleware())

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
			accountsGroup.GET("/search", accountsController.SearchAccounts)
			accountsGroup.POST("", accountsController.CreateAccount)
			accountsGroup.POST("/:id/withdraw", accountsController.Withdraw)
			accountsGroup.POST("/:id/deposit", accountsController.Deposit)

		}

		loansGroup := v1.Group("/loans")
		{
			loansGroup.POST("/calculate", loansController.AnnuitySchedule)
		}
	}

	r.Run(":" + os.Getenv("APP_PORT"))
}
