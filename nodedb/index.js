"use strict";

const mongodb = require("mongodb");
const MongoStore = require("./taskstore");
const express = require("express");
const app = express();

const addr = process.env.ADDR || "localhost:4000";
const [host, port] = addr.split(":");

const mongoAddr = process.env.DBADDR || "localhost:27017";
const mongoURL = `mongodb://${mongoAddr}/tasks`;

//connect to MongoDB
mongodb.MongoClient.connect(mongoURL)
    .then(db => {
        //initialize a new task store
        let taskStore = new MongoStore(db, "tasks");

        //parses posted JSON and makes
        //it available from req.body
        app.use(express.json());

        app.post("/v1/tasks", (req, res) => {
            //insert a new task
            //use taskStore.insert()
            let task = {
                title: req.body.title,
                completed: false
            }
            taskStore.insert(task)
                .then(task => {
                    res.json(task);
                })
                .catch(err => {
                    throw err;
                });
        });

        app.get("/v1/tasks", (req, res) => {
            //return all not-completed tasks in the database
            taskStore.getAll(false)
                .then(tasks => {
                    res.json(tasks);
                })
                .catch(err => {
                    throw err;
                });
        });

        app.patch("/v1/tasks/:taskID", (req, res) => {
            let taskIDToUpdate = new mongodb.ObjectID(req.params.taskID);
            let updates = {
                completed: req.body.completed
            };
            taskStore.update(taskIDToUpdate, updates)
                .then(updatedTask => {
                    res.json(updatedTask);
                })
                .catch(err => {
                    throw err;
                });
        });

        app.listen(port, host, () => {
            console.log(`server is listening at http://${addr}....`);
        });
    })
    .catch(err => {
        throw err;
    });