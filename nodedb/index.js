// @ts-check
"use strict";

const mongodb = require("mongodb");
const MongoStore = require("./taskstore");
const express = require("express");
const app = express();
const addr = process.env.ADDR || "localhost:4000";
const [host, port] = addr.split(":");
const portNum = parseInt(port);

const mongoAddr = process.env.DBADDR || "localhost:27017";
const mongoURL = `mongodb://${mongoAddr}/tasks`;

mongodb.MongoClient.connect(mongoURL)
    .then(db => {
        let taskStore = new MongoStore(db, "tasks");
        // parses posted json and makes
        // it available from req.body
        app.use(express.json());

        app.post("/v1/task", (req, res) => {
            // insert new task
            let task = {
                title: req.body.title,
                completed: false
            }
            taskStore.insert(task)
                .then(task => {
                    res.json(task);
                })
                .catch(err => {
                    throw err
                });
        });

        app.get("/v1/task", (req, res) => {
            // return all not-completed tasks in the database
            taskStore.getAll(false)
                .then(tasks => {
                    res.json(tasks);
                })
                .catch(err => {
                    throw err;
                });
        });

        app.patch("/v1/tasks/:taskID", (req, res) => {
            let taskIDToFetch = req.params.taskID;
            let updatedTask = {
                title: req.body.title,
                completed: req.body.completed
            }
            // fetch single task by id and send to client
            taskStore.update(taskIDToFetch, updatedTask)
                .then(task => {
                    res.json(task);
                })
                .catch(err => {
                    throw err;
                });
        });

        app.listen(portNum, host, () => {
            console.log(`server is listening at http://${addr}...`);
        });
    });