package main

import (
	scrapeops_plugin "github.com/Radicalius/scrapeops/shared"
)

type Context struct {
	Queue    *Queue
	Database *DatabaseCollection
	Logger   *Logger
}

func NewContext(q *Queue, db *DatabaseCollection) *Context {
	return &Context{
		Queue:    q,
		Database: db,
	}
}

func (c *Context) GetDatabase() scrapeops_plugin.Database {
	return c.Database
}

func (c *Context) GetQueue() scrapeops_plugin.Queue {
	return c.Queue
}

func (c *Context) GetLogger() scrapeops_plugin.Logger {
	return c.Logger
}

func (c *Context) WithLogger(logger *Logger) *Context {
	return &Context{
		Queue:    c.Queue,
		Database: c.Database,
		Logger:   logger,
	}
}
