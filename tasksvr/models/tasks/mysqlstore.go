package tasks

import (
	"database/sql"
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//various SQL statements we will need to execute
//TODO: write me some SQL! Refer to sql/schema.sql
//for the table/column names

//SQL to insert a new task row
//use `?` for column values that we will get at runtime
const sqlInsertTask = `insert into tasks(id,title,completed,createdAt) values (?,?,?,?)`

//SQL to insert a tag for a task
const sqlInsertTag = `insert into tags(taskID,tag) values (?,?)`

//SQL to select all tasks/tags with a particular task.completed value
//join tasks to tags so we get everything with only one query
const sqlSelectAll = `select id,title,completed,createdAt,modifiedAt,tag
from tasks inner join tags on (tasks.id=tags.taskID) where completed=?
order by id,createdAt`

//SQL to select a particular task/tags by id
//join tasks to tags so we get everything with only one query
const sqlSelectID = `select id,title,completed,createdAt,modifiedAt,tag
from tasks inner join tags on (tasks.id=tags.taskID) where id=?`

//SQL to update task.completed and task.modifiedAt where id = ?
const sqlUpdate = `update tasks set completed=?, modifiedAt=? where id=?`

//taskRow represents the data we will get for each row
//from the sqlSelectAll and selSelectID join queries
type taskRow struct {
	id         string
	title      string
	completed  bool
	createdAt  time.Time
	modifiedAt *time.Time
	tag        string
}

//MySQLStore implements Store for a MySQL database
type MySQLStore struct {
	//a live reference to the database
	db *sql.DB
}

//NewMySQLStore constructs a MySQLStore
func NewMySQLStore(db *sql.DB) *MySQLStore {
	if db == nil {
		panic("nil pointer passed to NewMySQLStore")
	}

	return &MySQLStore{
		db: db,
	}
}

//Insert inserts a new task into the database
func (s *MySQLStore) Insert(nt *NewTask) (*Task, error) {
	//convert the NewTask to a Task (which will also validate)
	task, err := nt.ToTask()
	if err != nil {
		return nil, err
	}
	//since we need to insert into both the `tasks` and `tags`
	//tables, and since we want those inserts to be atomic (all or nothing)
	//we need to start a database transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error begining transaction: %v", err)
	}

	//execute the insert to the `tasks` table
	//the .Hex() method of bson.ObjectId will return
	//the hexadecimal string representation of the binary
	//object ID, which is human-readable
	if _, err := tx.Exec(sqlInsertTask, task.ID.Hex(), task.Title, task.Completed, task.CreatedAt); err != nil {
		//rollback the transaction if there's an error
		tx.Rollback()
		return nil, fmt.Errorf("error inserting task: %v", err)
	}

	//for each tag in the Tags slice...
	for _, tag := range task.Tags {
		//...execute an insert to the `tags` table
		if _, err := tx.Exec(sqlInsertTag, task.ID.Hex(), tag); err != nil {
			//rollback the transaction if there's an error
			tx.Rollback()
			return nil, fmt.Errorf("error inserting tag for task: %v", err)
		}
	}

	//now commit the transaction so that all those inserts are atomic
	if err := tx.Commit(); err != nil {
		//try to rollback if we can't commit
		tx.Rollback()
		return nil, fmt.Errorf("error committing insert transaction: %v", err)
	}

	return task, nil
}

//GetAll gets all tasks with a particular completed state
func (s *MySQLStore) GetAll(completed bool) ([]*Task, error) {
	//execute the select all query, using the value of the
	//`completed` pararameter in the `where completed=?`
	//part of the SQL statement (the ? is replaced by the value)
	rows, err := s.db.Query(sqlSelectAll, completed)
	if err != nil {
		return nil, fmt.Errorf("error selecting tasks: %v", err)
	}

	//use our scanTasks function to scan the results into
	//a []*Task
	return scanTasks(rows)
}

//Update applies the values in `tu` to the task with `id`
func (s *MySQLStore) Update(id bson.ObjectId, tu *TaskUpdates) (*Task, error) {
	//execute the update SQL
	_, err := s.db.Exec(sqlUpdate, tu.Completed, time.Now(), id.Hex())
	if err != nil {
		return nil, fmt.Errorf("error updating task: %v", err)
	}

	//we need to return the updated Task, so select it
	//and scan the results
	rows, err := s.db.Query(sqlSelectID, id.Hex())
	if err != nil {
		return nil, fmt.Errorf("error selecting task after update: %v", err)
	}

	tasks, err := scanTasks(rows)
	if err != nil || len(tasks) == 0 {
		return nil, err
	}
	//return the first (and only)
	//element from the slice
	return tasks[0], nil
}

//scanTasks scans query result rows into a []*Task.
//Since the sqlSelectAll and sqlSelectID queries do a
//join between tasks and tags, we will get multiple rows
//for a single task if it has multiple tags. This function
//iterates over the result rows, constructing a *Task
//as it finds new task IDs, else it appends the next
//tag to the last-created Task.
func scanTasks(rows *sql.Rows) ([]*Task, error) {
	//ensure the rows are closed regardless of how
	//we exit this function
	defer rows.Close()
	//create an empty slice of *Task to hold
	//the results
	tasks := []*Task{}
	//create a taskRow struct to scan the data into
	row := taskRow{}
	//for each row in the resultset
	for rows.Next() {
		//scan the row data into our taskRow struct instance
		err := rows.Scan(&row.id, &row.title, &row.completed, &row.createdAt, &row.modifiedAt, &row.tag)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		//if this is not the first task, and if the current ID matches
		//the ID of the last task in the slice...
		if len(tasks) > 0 && tasks[len(tasks)-1].ID.Hex() == row.id {
			//...append just the tag to the slice of tags in that task
			lastTask := tasks[len(tasks)-1]
			lastTask.Tags = append(lastTask.Tags, row.tag)
		} else {
			//...create a new Task and append it to
			//the tasks slice
			task := &Task{
				ID:         bson.ObjectIdHex(row.id),
				Title:      row.title,
				Completed:  row.completed,
				CreatedAt:  row.createdAt,
				ModifiedAt: row.modifiedAt,
				Tags:       []string{row.tag},
			}
			tasks = append(tasks, task)
		}
	}

	//if there was an error reading rows off the network
	//rows.Err() will return a non-nil value
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %v", err)
	}

	return tasks, nil
}
