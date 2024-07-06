package shared

import (
	"github.com/jmoiron/sqlx"
)

type Context interface {
	GetQueue() Queue
	GetDatabase() Database
	GetLogger() Logger
	GetMetrics() Metrics
}

type Queue interface {
	Emit(funcName string, data interface{}) error
}

type Database interface {
	Exec(dbName string, sql string, params ...interface{}) error
	Query(dbName string, sql string, params ...interface{}) (*sqlx.Rows, error)
}

type Logger interface {
	Fatal(message string, params ...string)
	Error(message string, params ...string)
	Warn(message string, params ...string)
	Info(message string, params ...string)
}

type Metrics interface {
	IncrementCounter(counterName string)
	Observe(histName string, value float64)
}
