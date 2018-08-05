package handlers

import (
	"github.com/info344-a17/info344-in-class/tasksvr/models/tasks"
)

//Context holds context values
//used by multiple handler functions.
//see https://drstearns.github.io/tutorials/gohandlerctx/
type Context struct {
	//the tasks.Store to use for inserting,
	//getting, and updating tasks
	tasksStore tasks.Store
}

func NewHandlerContext(tasksStore tasks.Store) *Context {
	return &Context{
		tasksStore: tasksStore,
	}
}
