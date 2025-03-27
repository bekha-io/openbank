package types

type TransactionCategory string

const (
	TransactionCategoryDeposit       TransactionCategory = "deposit"
	TransactionCategoryWithdrawal    TransactionCategory = "withdrawal"
	TransactionCategoryTransfer      TransactionCategory = "transfer"
	TransactionCategoryTransport     TransactionCategory = "transport"
	TransactionCategoryGroceries     TransactionCategory = "groceries"
	TransactionCategoryEntertainment TransactionCategory = "entertainment"
	TransactionCategoryHealthcare    TransactionCategory = "healthcare"
	TransactionCategoryUtilities     TransactionCategory = "utilities"
	TransactionCategoryRent          TransactionCategory = "rent"
	TransactionCategorySalary        TransactionCategory = "salary"
	TransactionCategoryInvestment    TransactionCategory = "investment"
	TransactionCategoryInsurance     TransactionCategory = "insurance"
	TransactionCategoryDining        TransactionCategory = "dining"
	TransactionCategoryEducation     TransactionCategory = "education"
	TransactionCategoryShopping      TransactionCategory = "shopping"
	TransactionCategoryTax           TransactionCategory = "tax"
	TransactionCategoryCharity       TransactionCategory = "charity"
	TransactionCategoryTravel        TransactionCategory = "travel"
	TransactionCategoryFees          TransactionCategory = "fees"
)
