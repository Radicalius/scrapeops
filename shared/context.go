package shared

type Context interface {
	GetQueue() Queue
	GetDatabase() Database
}

type Queue interface {
	Emit(funcName string, data interface{}) error
}

type Database interface {
	Exec(sql string) error
	Query(sql string, out *[]interface{}) error
}
