package shared

func Query[T any](ctx Context, dbName string, sql string, out *[]T, params ...any) error {
	rows, err := ctx.GetDatabase().Query(dbName, sql, params...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var elem T
		rows.StructScan(&elem)
		*out = append(*out, elem)
	}

	return nil
}

func Exec(ctx Context, dbName string, sql string) error {
	return ctx.GetDatabase().Exec(dbName, sql)
}

func Emit[T any](ctx Context, queueName string, message T) error {
	return ctx.GetQueue().Emit(queueName, message)
}

func EmitHttp(ctx Context, callback string, url string, params ...string) error {
	req := HttpAsyncMessage{
		Callback: callback,
		Url:      url,
		Queries:  make([]Query_, 0),
	}

	for i := 0; i < len(params)/2; i++ {
		req.Queries = append(req.Queries, Query_{
			Selector:  params[i*2],
			Attribute: params[i*2+1],
		})
	}

	return ctx.GetQueue().Emit("httpAsync", req)
}
