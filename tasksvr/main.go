package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/info344-a17/info344-in-class/tasksvr/handlers"
	"github.com/info344-a17/info344-in-class/tasksvr/models/tasks"

	mgo "gopkg.in/mgo.v2"
)

const defaultAddr = ":80"

func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = defaultAddr
	}

	//MONGO DB CONNECTION
	//get the address of the MongoDB server
	//from an environment variable
	mongoAddr := os.Getenv("MONGO_ADDR")
	//default to "localhost"
	if len(mongoAddr) == 0 {
		mongoAddr = "localhost"
	}
	//dial the MongoDB server
	mongoSess, err := mgo.Dial(mongoAddr)
	if err != nil {
		log.Fatalf("error dialing mongo: %v", err)
	}
	//construct a new MongoStore, providing the mgo.Session
	//as well as a database and collection name to use
	mongoStore := tasks.NewMongoStore(mongoSess, "tasks", "tasks")

	//MYSQL CONNECTION
	//get the address for the MySQL server
	mysqlAddr := os.Getenv("MYSQL_ADDR")
	if len(mysqlAddr) == 0 {
		mysqlAddr = "localhost"
	}
	//construct the connection string
	mysqlConfig := mysql.NewConfig()
	mysqlConfig.Addr = mysqlAddr
	mysqlConfig.DBName = os.Getenv("MYSQL_DATABASE")
	mysqlConfig.User = "root"
	mysqlConfig.Passwd = os.Getenv("MYSQL_ROOT_PASSWORD")
	//tell the MySQL driver to parse DATETIME
	//column values into go time.Time values
	mysqlConfig.ParseTime = true

	db, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		log.Fatalf("error opening mysql database: %v", err)
	}
	defer db.Close()
	//un-comment this next line to create a MySQLStore
	//mysqlStore := tasks.NewMySQLStore(db)

	//construct a new handler context passing the MongoStore
	//as the tasks.Store implementation to use. The handlers
	//only work with the abstract tasks.Store interface so they
	//don't have to care which store is being used
	//see https://drstearns.github.io/tutorials/gohandlerctx/
	handlerCtx := handlers.NewHandlerContext(mongoStore)

	//un-comment this next line to use MySQL instead of
	//MongoDB as the tasks.Store implementation
	//handlerCtx = handlers.NewHandlerContext(mysqlStore)

	mux := http.NewServeMux()
	//because TasksHandler and SpecificTasksHandler are methods
	//of the handlers.Context struct, the will have access to
	//the fields of that struct when they are called by the mux.
	//see https://drstearns.github.io/tutorials/gohandlerctx/
	mux.HandleFunc("/v1/tasks", handlerCtx.TasksHandler)
	mux.HandleFunc("/v1/tasks/", handlerCtx.SpecificTaskHandler)

	fmt.Printf("server is listening at http://%s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
