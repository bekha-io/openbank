package main

import (
	"context"
	"os"

	"github.com/bekha-io/openbank/domain/entities/permissions"
	"github.com/bekha-io/openbank/domain/services"
	"github.com/bekha-io/openbank/infrastructure/repository/file"
	"github.com/bekha-io/openbank/infrastructure/repository/mongodb"
	"github.com/bekha-io/openbank/presentation/rest/accounts"
	"github.com/bekha-io/openbank/presentation/rest/auth"
	"github.com/bekha-io/openbank/presentation/rest/customers"
	"github.com/bekha-io/openbank/presentation/rest/employees"
	"github.com/bekha-io/openbank/presentation/rest/loans"
	"github.com/bekha-io/openbank/presentation/rest/me"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
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

	dbName := "openbank"
	transactionsRepo := mongodb.NewMongoTransactionRepository(mongoClient, dbName)
	accountsRepo := mongodb.NewMongoAccountRepository(mongoClient, dbName)
	individualCustomersRepo := mongodb.NewMongoIndividualCustomerRepository(mongoClient, dbName)
	employeesRepo := mongodb.NewMongoEmployeeRepository(mongoClient, dbName)
	loansRepo := mongodb.NewMongoLoanRepository(mongoClient, dbName)
	rolesRepo := &file.FileRolesRepository{Dir: "tmp/roles"}

	accountsSvc := services.NewAccountsService(accountsRepo, transactionsRepo)
	individualCustomersSvc := services.NewIndividualCustomerService(individualCustomersRepo, accountsRepo)
	employeesSvc := services.NewEmployeeService(employeesRepo)
	loansSvc := services.NewLoanService(loansRepo)
	authSvc := services.NewAuthorizationService(rolesRepo)

	accountsController := accounts.NewAccountsController(accountsSvc)
	customersController := customers.NewCustomerController(individualCustomersSvc)
	employeesController := employees.NewEmployeeController(employeesSvc)
	loansController := loans.NewLoanController(loansSvc)
	authController := auth.AuthController{EmployeeService: employeesSvc, AuthorizationService: authSvc}

	r := gin.Default()
	r.Use(CORSMiddleware())

	// api/v1
	v1 := r.Group("/v1")
	v1.POST("/employees/auth", authController.EmployeeSignIn)
	v1.Use(authController.EmployeeAuthenticateMiddleware())
	{
		// api/v1/customers
		customersGroup := v1.Group("/customers")
		{
			customersGroup.GET("/:id", customersController.GetCustomer)
			customersGroup.GET("/search", customersController.SearchCustomers)
			customersGroup.POST("", customersController.CreateCustomer)
			customersGroup.GET("/:id/accounts", customersController.GetCustomerAccounts)
		}

		accountsGroup := v1.Group("/accounts")
		{
			accountsGroup.GET("/search", accountsController.SearchAccounts)
			accountsGroup.POST("/transfer", accountsController.Transfer)
			accountsGroup.GET("/:id", accountsController.GetAccount)
			accountsGroup.GET("/:id/transactions", accountsController.GetAccountTransactions)
			accountsGroup.POST("", accountsController.CreateAccount)
			accountsGroup.POST("/:id/withdraw", accountsController.Withdraw)
			accountsGroup.POST("/:id/deposit", accountsController.Deposit)

		}

		employeesGroup := v1.Group("/employees")
		{
			employeesGroup.GET("/search", authController.IfHasPermissions(employeesController.SearchEmployees, permissions.CanReadEmployee))
			employeesGroup.POST("/", employeesController.CreateEmployee)
		}

		loansGroup := v1.Group("/loans")
		{
			loansGroup.POST("/calculate", loansController.AnnuitySchedule)

			// Loan products
			loansGroup.GET("/products", loansController.GetLoanProducts)
			loansGroup.POST("/products", loansController.CreateLoanProduct)
		}
	}

	meCtrl := me.NewController(accountsSvc, individualCustomersSvc)

	// Temporary API for checking mobile application.
	// Mobile app should be a separate solution, not attached to core
	// /me
	meg := r.Group("/me")
	meg.POST("/auth", meCtrl.CustomerSignIn)
	meg.Use(meCtrl.CustomerAuthenticateMiddleware())
	{
		meg.GET("/accounts", meCtrl.GetAccounts)
	}

	r.Run(":" + os.Getenv("APP_PORT"))
}
