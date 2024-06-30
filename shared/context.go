package shared

import (
	"github.com/jmoiron/sqlx"
)

type Context interface {
	GetQueue() Queue
	GetDatabase() Database
}

type Queue interface {
	Emit(funcName string, data interface{}) error
}

type Database interface {
	Exec(dbName string, sql string, params ...interface{}) error
	Query(dbName string, sql string, params ...interface{}) (*sqlx.Rows, error)
}
