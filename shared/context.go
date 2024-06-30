package shared

import "reflect"

type Context interface {
	GetQueue() Queue
	GetDatabase() Database
}

type Queue interface {
	Emit(funcName string, data interface{}) error
}

type Database interface {
	Exec(dbName string, sql string, params ...interface{}) error
	Query(t reflect.Type, dbName string, sql string, out *[]interface{}, params ...interface{}) error
}
