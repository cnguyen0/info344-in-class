// @ts-check
// ^ this does type check in visual studio code.
"use strict";

// load the express and morgan modules
const express = require("express");
const morgan = require("morgan");
const handlers = require("./handlers");

const addr = process.env.ADDR || ":80"; // do the || to default
const [host, port] = addr.split(":");
const portNum = parseInt(port);

// creates an express app for us to use. we can now use it as an app.
const app = express();
app.use(morgan(process.env.LOG_FORMAT || "dev")); // app is an object.
// use will execute every single request. similar to a handler.
// express is different than go mux. we have to wrap the mux with a middleware.
// express, you an add multiple of this guy and it put it into a chain of requests
// morgan will have the chance to preprocess the request
// morgan = middleware handler. and you can add multiple by doing app.use()
// accepts a logging format! ez. // "dev" was hardcoded. but you can use anything tbh bc u want
// to let the people starting the server decide.

// "/" root directory
// the second parameter is a function lmaooooooooo
app.get("/", (req, res) => {
    res.set("Content-Type", "text/plain"); // default content type is text/html
    res.send("Hello, Node.js!");
});

app.get("/v1/users/me/hello", (req, res) => {
    // this gateway will call this server and will send a specialized header
    let userJSON = req.get("X-User"); // "X-" is a prefix for nonstandard header
    if (!userJSON) {
        // this throws it into the html so you need to hide it on the html and
        // put it in the terminal only!!!!!
        // see the last app.use() below before the app.listen()
        throw new Error("No X-User header provided");
    }
    let user = JSON.parse(userJSON);
    res.json({
        // the backtick will allow u to include javascript variable in the string
        message: `Hello, ${user.firstName} ${user.lastName}`
    });
});

app.use(handlers({}));

// last handler you need to put so it will prevent the error from showing in the html
// protects from hacker
app.use((err, req, res, next) => {
    console.error(err.stack); // that is the full stack of the err that was thrown
    res.set("Content-Type", "text/plain");
    res.send(err.message); // send the message only to the client
});

// last parameter is a lamda function that will call when its up and running
app.listen(portNum, host, () => {
    console.log(`server is listening at http://${addr}...`);
});