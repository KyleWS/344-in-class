package tasks

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//completedFilter is used to filter the set
//of tasks we find during the GetAll() method
type completedFilter struct {
	Completed bool
}

//updateDoc is used to update fields in an existing Task
type updateDoc struct {
	Completed  bool
	ModifiedAt time.Time
}

//MongoStore implements Store for MongoDB
type MongoStore struct {
	//the mongo session
	session *mgo.Session
	//the database name to use
	dbname string
	//the collection name to use
	colname string
	//the Collection object for that dbname/colname
	col *mgo.Collection
}

//NewMongoStore constructs a new MongoStore, given a live mgo.Session,
//a database name, and a collection name
func NewMongoStore(sess *mgo.Session, dbName string, collectionName string) *MongoStore {
	if sess == nil {
		//panic will halt the program and print
		//a stack trace to the terminal so that
		//developers can quickly find where the
		//panic occurred
		//(log.Fatal() doesn't print a stack trace)
		panic("nil pointer passed for session")
	}

	//return a new MongoStore
	return &MongoStore{
		session: sess,
		dbname:  dbName,
		colname: collectionName,
		//store a reference to the mgo.Collection
		//for the named database/collection so that
		//we don't have to get it during each method
		col: sess.DB(dbName).C(collectionName),
	}
}

//Insert inserts a new task into the store
func (s *MongoStore) Insert(nt *NewTask) (*Task, error) {
	//convert the NewTask to a Task,
	//which will also validate the NewTask
	task, err := nt.ToTask()
	if err != nil {
		return nil, err
	}

	//insert the Task into the database and return it,
	//or an error if one occurred
	if err := s.col.Insert(task); err != nil {
		return nil, fmt.Errorf("error inserting task: %v", err)
	}
	return task, nil
}

//GetAll gets all tasks (up to AllTasksLimit) with a particular `completed` value
func (s *MongoStore) GetAll(completed bool) ([]*Task, error) {
	//create an empty slice of *Task to hold the results
	tasks := []*Task{}

	//create a filter document to get only the tasks where
	//the Completed field is set to the value of the
	//`completed` parameter
	filter := &completedFilter{completed}

	//find tasks that satisfy the filter, but limit
	//the number of results to our AllTasksLimit.
	//the .All() method will populate the `tasks` slice
	//with all of the Tasks returned from the database
	if err := s.col.Find(filter).Limit(AllTasksLimit).All(&tasks); err != nil {
		return nil, fmt.Errorf("error getting tasks: %v", err)
	}
	return tasks, nil
}

//Update updates the task with `id` based on the values in `tu`
func (s *MongoStore) Update(id bson.ObjectId, tu *TaskUpdates) (*Task, error) {
	//create an update document, which lists the fields to update
	//and the values to update them to
	upd := &updateDoc{
		Completed:  tu.Completed,
		ModifiedAt: time.Now(),
	}
	//create a change document, which specifies the updates to
	//make, plus any options for the update operation. Here we
	//set ReturnNew to true, so that we get the updated version
	//of the Task (default behavior is to return the previous
	//un-updated version of the Task)
	change := mgo.Change{
		Update:    bson.M{"$set": upd},
		ReturnNew: true,
	}

	//create an empty Task instance to receive the updated Task
	task := &Task{}

	//find the task to update using the provided `id`
	//and apply the change document to it. `task` will be
	//populated with the updated Task data.
	if _, err := s.col.FindId(id).Apply(change, task); err != nil {
		return nil, fmt.Errorf("error updating task: %v", err)
	}
	return task, nil
}
