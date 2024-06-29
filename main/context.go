package main

import (
	scrapeops_plugin "github.com/Radicalius/scrapeops/shared"
)

type Context struct {
	Queue    *Queue
	Database *Database
}

func NewContext(q *Queue, db *Database) *Context {
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
