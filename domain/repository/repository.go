package repository


type Repository struct {
	IndividualCustomers IIndividualCustomerRepository
	Accounts IAccountRepository
}