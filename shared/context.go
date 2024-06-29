package shared

type Context interface {
	GetQueue() Queue
	GetDatabase() Database
}

type Queue interface {
	Emit(funcName string, data interface{}) error
}

type Database interface {
	Exec(dbName string, sql string) error
	Query(dbName string, sql string, out *[]interface{}) error
}
