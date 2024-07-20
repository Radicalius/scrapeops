package main

import (
	scrapeops_plugin "github.com/Radicalius/scrapeops/shared"
	"gorm.io/gorm"
)

type Context struct {
	Queue              *Queue
	DatabaseCollection *DatabaseCollection
	Metrics            *Metrics
	Logger             *Logger
}

func NewContext(q *Queue, db *DatabaseCollection, m *Metrics) *Context {
	return &Context{
		Queue:              q,
		DatabaseCollection: db,
		Metrics:            m,
	}
}

func (c *Context) GetDatabase(dbName string) *gorm.DB {
	return c.DatabaseCollection.GetDatabase(dbName)
}

func (c *Context) GetQueue() scrapeops_plugin.Queue {
	return c.Queue
}

func (c *Context) GetLogger() scrapeops_plugin.Logger {
	return c.Logger
}

func (c *Context) GetMetrics() scrapeops_plugin.Metrics {
	return c.Metrics
}

func (c *Context) WithLogger(logger *Logger) *Context {
	return &Context{
		Queue:              c.Queue,
		DatabaseCollection: c.DatabaseCollection,
		Logger:             logger,
	}
}
