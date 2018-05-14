"use strict";

//SQL statements
const SQL_SELECT_ALL_CHANNELS = "select * from channels c";
const SQL_SELECT_ALL_CHANNELS_FOR_USER = SQL_SELECT_ALL_CHANNELS +
                                            " join channel_user cu on cu.channelID = c.id" +
                                            " join users u on u.id = cu.userID" +
                                            " where (cu.userID = ? or c.private = 'false')" +
                                            " order by c.id";

const SQL_SELECT_BY_ID = SQL_SELECT_ALL + " where id=?";
const SQL_INSERT = "insert into links (url,comment) values (?,?)";
const SQL_UPVOTE = "update links set votes=votes+1 where id=?";

const express = require("express");
const mysql = require("mysql");


const app = express();

let db = mysql.createPool({
    host: process.env.MYSQL_ADDR,
    database: process.env.MYSQL_DATABASE,
    user: "root",
    password: process.env.MYSQL_ROOT_PASSWORD
});

//checkUserAuth
function checkUserAuth(req) {
    let userJson = req.get("X-User");
    if (!userJson) {
        return new Error("Please Sign In")
    }
}

app.use(express.json());

// Handle Endpoint: /v1/channels
app.get("/v1/channels", (req, res, next) => {
    var err = checkUserAuth(req);
    if (err) {
        return next(err);
    }
    let channels = [];

    db.Query(SQL_SELECT_ALL_CHANNELS, (err, rows) => {
        if (err) {
            return next(err)
        }
        rows.forEach((row) => {

        })
    })
});

app.post("/v1/channels", (req, res, next) => {
    var err = checkUserAuth(req);
    if (err) {
        return next(err);
    }
});

// Handle Endpoint: /v1/channels/{channelID}
app.get("/v1/channels/:channelID", (req, res, next) => {
    var err = checkUserAuth(req);
    if (err) {
        return next(err);
    }
});

app.post("/v1/channels/:channelID", (req, res, next) => {
    var err = checkUserAuth(req);
    if (err) {
        return next(err);
    }
});

app.patch("/v1/channels/:channelID", (req, res, next) => {
    var err = checkUserAuth(req);
    if (err) {
        return next(err);
    }
});

app.delete("/v1/channels/:channelID", (req, res, next) => {
    var err = checkUserAuth(req);
    if (err) {
        return next(err);
    }
});

// Handle Endpoint: /v1/channels/{channelID}/members
app.post("/v1/channels/:channelID/members", (req, res, next) => {
    var err = checkUserAuth(req);
    if (err) {
        return next(err);
    }
});

app.delete("/v1/channels/:channelID/members", (req, res, next) => {
    var err = checkUserAuth(req);
    if (err) {
        return next(err);
    }
});

// Handle Endpoint: /v1/messages/{messageID}
app.patch("/v1/messages/:messageID", (req, res, next) => {
    var err = checkUserAuth(req);
    if (err) {
        return next(err);
    }
});

app.delete("/v1/messages/:messageID", (req, res, next) => {
    var err = checkUserAuth(req);
    if (err) {
        return next(err);
    }
});

app.use((err, req, res, next) => {
    if (err.stack) {
        console.error(err.stack);
    }
    res.status(500).send("Something bad happened. Sorry.");
});