package shared

import "gorm.io/gorm"

type Context interface {
	GetQueue() Queue
	GetDatabase(dbName string) *gorm.DB
	GetLogger() Logger
	GetMetrics() Metrics
}

type Queue interface {
	Emit(funcName string, data interface{}) error
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
