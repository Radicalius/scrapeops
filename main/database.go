package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	conn *sqlx.DB
}

func InitDatabase(dbName string, migrations []string) (*Database, error) {
	path := os.Getenv("SCRAPEOPS_DATABASE_DIRECTORY")
	if path == "" {
		path = "./dbs"
	}

	filePath := fmt.Sprintf("%s/%s.db", path, dbName)

	db, err := sqlx.Connect("sqlite3", filePath)
	if err != nil {
		return nil, err
	}

	for _, migration := range migrations {
		_, err := db.Exec(migration)
		if err != nil {
			return nil, fmt.Errorf("Error applying migration: \n\tdatabase: %s\n\tstep: %s\n\terror: %s", dbName, err.Error())
		}
	}

	return &Database{
		conn: db,
	}, nil
}

func (d *Database) Exec(sql string) error {
	_, err := d.conn.Exec(sql)
	return err
}

func (d *Database) Query(sql string, out *[]interface{}) error {
	rows, err := d.conn.Queryx(sql)
	if err != nil {
		return err
	}

	elemType := reflect.TypeOf(out).Elem()
	for rows.Next() {
		newElem := reflect.New(elemType).Interface()
		err = rows.StructScan(newElem)
		if err != nil {
			return err
		}

		*out = append(*out, newElem)
	}

	return nil
}

type DatabaseCollection struct {
	dbs map[string]*Database
}

func NewDatabaseCollection() *DatabaseCollection {
	return &DatabaseCollection{
		dbs: make(map[string]*Database),
	}
}

func (dbc *DatabaseCollection) AddDatabase(dbName string, migrations []string) error {
	db, err := InitDatabase(dbName, migrations)
	if err != nil {
		return err
	}

	dbc.dbs[dbName] = db
	return nil
}

func (dbc *DatabaseCollection) Exec(dbName string, sql string) error {
	db, exists := dbc.dbs[dbName]
	if !exists {
		return fmt.Errorf("Database %s does not exist", dbName)
	}

	return db.Exec(sql)
}

func (dbc *DatabaseCollection) Query(dbName string, sql string, out *[]interface{}) error {
	db, exists := dbc.dbs[dbName]
	if !exists {
		return fmt.Errorf("Database %s does not exist", dbName)
	}

	return db.Query(sql, out)
}
