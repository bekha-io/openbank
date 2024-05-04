package repository

type Filter struct {
	Key string

	EqualTo interface{}
	GreaterThan interface{}
	LessThan interface{}
	Like interface{}
	Before interface{}
	After interface{}
}