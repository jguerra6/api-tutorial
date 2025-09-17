package postgres

import (
	"errors"
)

var (
	ErrSQLBegin     = errors.New("unable to start a transaction")
	ErrRollback     = errors.New("unable to rollback the transaction")
	ErrSQLExec      = errors.New("sql query exec error")
	ErrDataNotFound = errors.New("data not found")
)
