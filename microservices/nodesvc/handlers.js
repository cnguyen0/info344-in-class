// @ts-check
"use strict";

const express = require("express");

// node modules don't do like export stuff
// give the module.export whatever the fuck u want
module.exports = (mongoSession) => {
    if (!mongoSession) {
        throw new Error("provide Mongo Session");
    }

    let router = express.Router();
    router.get("/v1/channels", (req, res) => {
        // we will talk to our mongo db in hurrrrr
        // query mongo using mongo session (will demonstrate on tuesday...)
        res.json([{name: "general"}]); // hardcode some shit in dw lol
    });

    return router;
}