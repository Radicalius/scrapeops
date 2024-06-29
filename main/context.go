package main

import (
	scrapeops_plugin "github.com/Radicalius/scrapeops/shared"
)

type Context struct {
	Queue    *Queue
	Database *DatabaseCollection
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
