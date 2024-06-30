package shared

import (
	"fmt"
	"reflect"
)

func Query[T any](ctx Context, dbName string, sql string, out *[]T, params ...any) error {
	t := reflect.TypeOf(*out).Elem()

	buffer := make([]interface{}, 0)
	fmt.Println(t.Name())
	err := ctx.GetDatabase().Query(t, dbName, sql, &buffer, params...)
	if err != nil {
		return err
	}

	for _, item := range buffer {
		*out = append(*out, item.(T))
	}
	
	return nil
}
