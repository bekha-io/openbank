package permissions

var (
	CanCrudLoanProduct   = NewPermission("loan_products", CRUD)
	CanCreateLoanProduct = NewPermission("loan_products", Create)
	CanUpdateLoanProduct = NewPermission("loan_products", Update)
	CanDeleteLoanProduct = NewPermission("loan_products", Delete)
	CanReadLoanProduct   = NewPermission("loan_products", Read)
)

var LoanProductPermissionSet = []Permission{
	CanCrudLoanProduct, CanCreateLoanProduct, CanUpdateLoanProduct,
	CanReadLoanProduct, CanDeleteLoanProduct,
}
