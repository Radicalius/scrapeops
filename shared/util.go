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
