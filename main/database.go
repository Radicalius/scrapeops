package main

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabase(dbName string, tables []interface{}) (*gorm.DB, error) {
	path := os.Getenv("SCRAPEOPS_DATABASE_DIRECTORY")
	if path == "" {
		path = "./dbs"
	}

	filePath := fmt.Sprintf("%s/%s.db", path, dbName)

	db, err := gorm.Open(sqlite.Open(filePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(tables...)

	return db, nil
}

type DatabaseCollection struct {
	dbs map[string]*gorm.DB
}

func NewDatabaseCollection() *DatabaseCollection {
	return &DatabaseCollection{
		dbs: make(map[string]*gorm.DB),
	}
}

func (dbc *DatabaseCollection) AddDatabase(dbName string, tables []interface{}) error {
	db, err := InitDatabase(dbName, tables)
	if err != nil {
		return err
	}

	dbc.dbs[dbName] = db
	return nil
}

func (dbc *DatabaseCollection) GetDatabase(dbName string) *gorm.DB {
	db, exists := dbc.dbs[dbName]
	if !exists {
		return nil
	}

	return db
}
