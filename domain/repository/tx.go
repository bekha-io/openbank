package repository

type Tx interface {
	Commit() error
	Rollback() error
}