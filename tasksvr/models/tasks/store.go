package tasks

import (
	"gopkg.in/mgo.v2/bson"
)

//AllTasksLimit is the maximum number of tasks
//a store should fetch, so that we don't explode
//when there are millions of tasks in the DBMS.
//We can add result set paging later.
const AllTasksLimit = 1000

//Store represents a data store for Tasks.
//The HTTP handler functions will use this abstract
//interface so that they don't have to care which
//DBMS we actually use. In this demo there are two
//implementations: mongostore.go, which is backed by
//MongoDB; and mysqlstore.go, which is backed by MySQL
type Store interface {
	//Insert inserts a NewTask into the database and returns
	//the inserted Task or an error
	Insert(nt *NewTask) (*Task, error)
	//GetAll selects all tasks (up to AllTasksLimit) where
	//the Completed field equals the completed parameter value
	GetAll(completed bool) ([]*Task, error)
	//Update updates the task matching the supplied `id`,
	//setting the fields in `tu` to the supplied values
	Update(id bson.ObjectId, tu *TaskUpdates) (*Task, error)
}
