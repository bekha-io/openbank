package main

import (
	"os"

	"github.com/bekha-io/openbank/domain/services"
	"github.com/bekha-io/openbank/infrastructure/fineract"
	"github.com/bekha-io/openbank/infrastructure/repository/gorm"
	"github.com/bekha-io/openbank/presentation/rest/me"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/shopspring/decimal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	fineractBaseUrl := os.Getenv("FINERACT_BASE_URL")
	fineractUsername := os.Getenv("FINERACT_USERNAME")
	fineractTenant := os.Getenv("FINERACT_TENANT")
	fineractPassword := os.Getenv("FINERACT_PASSWORD")

	fr := &fineract.FineractClient{
		BaseUrl:    fineractBaseUrl,
		TenantName: fineractTenant,
		Username:   fineractUsername,
		Password:   fineractPassword,
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DB_URI")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	gormRepo.Init(db)
	beneficiariesMemoryRepo := gormRepo.NewPostgresBenificiaryRepository(db)
	transactionsMemoryRepo := gormRepo.NewPostgresTransactionRepository(db)

	accountsSvc := services.NewAccountsService(fr, beneficiariesMemoryRepo, transactionsMemoryRepo)
	individualCustomersSvc := services.NewIndividualCustomerService(fr, beneficiariesMemoryRepo)

	r := gin.Default()
	r.Use(CORSMiddleware())

	meCtrl := me.NewController(accountsSvc, individualCustomersSvc)

	// Temporary API for checking mobile application.
	// Mobile app should be a separate solution, not attached to core
	// /me
	meg := r.Group("/me")
	meg.POST("/auth", meCtrl.CustomerSignIn)
	meg.Use(meCtrl.CustomerAuthenticateMiddleware())
	{
		meg.GET("/accounts", meCtrl.GetAccounts)
		meg.GET("/accounts/:id/transactions", meCtrl.GetAccountTransactions)
		meg.GET("/customers/:phoneNumber", meCtrl.GetCustomerByPhoneNumber)

		bnf := meg.Group("/beneficiaries")
		{
			bnf.POST("", meCtrl.CreateBeneficiary)
			bnf.GET("", meCtrl.GetCustomerBeneficiaries)
			bnf.GET("/:id", meCtrl.GetBeneficiaryByID)
		}

		transfers := meg.Group("/transfers")
		{
			transfers.POST("", meCtrl.TransferMoney)
		}
	}

	r.Run(":" + os.Getenv("APP_PORT"))
}
