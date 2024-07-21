package shared

import "gorm.io/gorm"

func Emit[T any](ctx Context, queueName string, message T) error {
	return ctx.GetQueue().Emit(queueName, message)
}

func EmitHttp(ctx Context, callback string, url string, joinKey string, params ...string) error {
	req := HttpAsyncMessage{
		Callback: callback,
		Url:      url,
		Queries:  make([]Query_, 0),
		JoinKey:  joinKey,
	}

	for i := 0; i < len(params)/2; i++ {
		req.Queries = append(req.Queries, Query_{
			Selector:  params[i*2],
			Attribute: params[i*2+1],
		})
	}

	return ctx.GetQueue().Emit("httpAsync", req)
}

func GetDatabase(ctx Context) *gorm.DB {
	return ctx.GetDatabase(DatabaseName)
}
