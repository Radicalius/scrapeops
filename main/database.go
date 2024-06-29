package main

type Database struct {
}

func (d *Database) Exec(sql string) error {
	return nil
}

func (d *Database) Query(sql string, out *[]interface{}) error {
	return nil
}
