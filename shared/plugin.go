package shared

import "encoding/json"

type RawHandlerFunc func([]byte, *Context) error
type HandlerFunc[T interface{}] func(T, *Context) error
type RawHandlerFuncMap map[string]RawHandlerFunc

var Handlers RawHandlerFuncMap = make(RawHandlerFuncMap)

func Register[T interface{}](name string, callback HandlerFunc[T]) {
	Handlers[name] = func(inp []byte, ctx *Context) error {
		var message T
		err := json.Unmarshal(inp, &message)
		if err != nil {
			return err
		}

		return callback(message, ctx)
	}
}
