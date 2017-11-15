#!/usr/bin/env node
"use strict"

const qName = "testQ";
const mqAddr = process.env.MQADDR || "localhost:5672";
const mqURL = `amqp://${mqAddr}`;
const amqp = require("amqplib");

// immediate execute
(async function() {
    console.log("connecting to %s", mqURL);
    let connection = await amqp.connect(mqURL);
    let channel = await connection.createChannel(); // await keyword will turn it into a syncronous call!!!1
    let qConf = await channel.assertQueue(qName, {durable: false});

    let testObj = {
        firstName: "Cindy",
        lastName: "Nguyen"
    }

    // testJSON = JSON.stringify(testObj);

    console.log("starting to send messages....");
    setInterval(() => {
        let msg = "Message sent at " + new Date().toLocaleTimeString();
        let testJSON = JSON.stringify(testObj);
        channel.sendToQueue(qName, Buffer.from(testJSON));
    }, 1000);
})();