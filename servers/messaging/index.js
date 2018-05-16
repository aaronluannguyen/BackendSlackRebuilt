"use strict";

let Channel = require("./models/channel");
let Message = require("./models/message");

//SQL statements
const SQL_SELECT_CHANNEL_BY_ID = "select * from channels where id=?";
const SQL_UPDATE_CHANNEL_NAME_DESC = "update channels set name=? set desc=? where id=?";
const SQL_DELETE_CHANNEL_AND_MESSAGES = "delete from c, cu, m" +
                                        " from channels c" +
                                        " join messages m on m.channelID = c.id" +
                                        " join channel_user cu on cu.channelID = c.id" +
                                        " where cu.channelID=?";
const SQL_SELECT_ALL_CHANNELS_FOR_USER = "select * from channels c" +
                                            " join channel_user cu on cu.channelID = c.id" +
                                            " join users u on u.id = cu.userID" +
                                            " where (c.private=? or cu.userID=?)" +
                                            " order by c.id";
const SQL_INSERT_NEW_CHANNEL = "insert into channels (name, description, private, createdAt, creatorUserID, editedAt) " +
                                "values (?,?,?,?,?,?)";
const SQL_INSERT_INTO_CHANNEL_USER = "insert into channel_user (userID, channelID) values(?,?)";
const SQL_DELETE_USER_FROM_CHANNEL = "delete from channel_user where (channelID=? and userID=?)";
const SQL_SELECT_SPECIFIC_CHANNEL = "select * from channels c" +
                                    " join channel_user cu on cu.channelID = c.id" +
                                    " where c.id = ?";
const SQL_TOP_100_MESSAGES = "select * from channels c" +
                                " join messages m on m.channelID = c.id" +
                                " where c.id = ?" +
                                " order by m.createdAt desc" +
                                " limit 100";

const SQL_POST_MESSAGE = "insert into messages (channelID, body, createdAt, creatorUserID, editedAt) values (?,?,?,?,?)";
const SQL_SELECT_MESSAGE_BY_ID = "select * from messages where id=?";
const SQL_UPDATE_MESSAGE= "update messages set body=? set editedAt=? where id=?";
const SQL_DELETE_MESSAGE_BY_ID = "delete from messages where id?";

const express = require("express");
const mysql = require("mysql");


const app = express();

const addr = process.env.ADDR || ":80";
const [host, port] = addr.split(":");

let db = mysql.createPool({
    host: process.env.MYSQL_ADDR,
    database: process.env.MYSQL_DATABASE,
    user: "root",
    password: process.env.MYSQL_ROOT_PASSWORD
});

//checkUserAuth checks the x-user header to make sure the user is signed in
//if user is signed in, the user, as json, is returned
//if not signed in, returns false
function checkUserAuth(req) {
    let userJson = req.get("X-User");
    if (!userJson) {
        return false;
    }
    let user = JSON.parse(userJson);
    return user;
}

//verifyUserInChannel makes sure the user is a member of a private channel
//and returns true or false accordingly
function verifyUserInChannel(userID, channelID, next) {
    db.query(SQL_SELECT_SPECIFIC_CHANNEL, [channelID], (err, rows) => {
        if (err) {
            return next(err);
        }
        if (rows[0].private === false) {
            return true;
        }
        rows.forEach((row) => {
           if (row.userID === userID) {
               return true;
           }
        });
        return false;
    });
}

//checkUserIsCreator checks if the user is the creator of a specified channel
//and returns true or false accordingly
function checkUserIsCreator(userID, next) {
    db.query(SQL_SELECT_CHANNEL_BY_ID, [userID], (err, rows) => {
        if (err) {
            return next(err);
        }
        return rows[0].creatorUserID === userID;
    });
}

//checkUserIsMessageCreator checks if the user is the creator of a specified message
//and returns true or false accordingly
function checkUserIsMessageCreator(userID, messageID, next) {
    db.query(SQL_SELECT_MESSAGE_BY_ID, [messageID], (err, rows) => {
        if (err) {
            return next(err);
        }
        return row[0].creatorUserID === userID;
    });
}

app.use(express.json());

// Handle Endpoint: /v1/channels
app.get("/v1/channels", (req, res, next) => {
    let user = checkUserAuth(req);
    if (user === false) {
        return res.status(401).send("Please sign in");
    }

    let channels = [];

    db.query(SQL_SELECT_ALL_CHANNELS_FOR_USER, [false, user.id], (err, rows) => {
        if (err) {
            return next(err);
        }
        console.log(rows.length);
        let members = [];
        let creator = {id: rows[0].userID, userName: rows[0].username,
                        firstName: rows[0].firstName, lastName: rows[0].lastName, photoURL: rows[0].photoURL};
        let currChannelID = rows[0].channelID;
        let channel = new Channel(rows[0].channelID, rows[0].name, rows[0].description, rows[0].private,
                                    members, rows[0].createdAt, creator, rows[0].editedAt);

        rows.forEach((row) => {
            if (row.channelID !== currChannelID) {
                channels.push(channel);
                members = [row.userID];
                creator = {id: row.userID, userName: row.username,
                    firstName: row.firstName, lastName: row.lastName, photoURL: row.photoURL};
                channel = new Channel(row.channelID, row.name, row.description, row.private, members,
                                        row.createdAt, creator, row.editedAt);
                currChannelID = row.id;
            } else {
                members.push(row.userID);
                if (row.creatorUserID === row.userID) {
                    creator = {id: row.userID, userName: row.username,
                        firstName: row.firstName, lastName: row.lastName, photoURL: row.photoURL};
                    channel.creator = creator;
                }
            }
        });
        channels.push(channel);
        res.json(channels);
    });
});

app.post("/v1/channels", (req, res, next) => {
    let user = checkUserAuth(req);
    if (user === false) {
        return res.status(401).send("Please sign in");
    }

    if (!req.body.name) {
        return res.status(400).send("Please provide name for channel");
    }
    let currTimestamp = new Date().toISOString().slice(0, 19).replace('T', ' ');
    db.query(SQL_INSERT_NEW_CHANNEL, [req.body.name, req.body.description, req.body.private,
                currTimestamp, user.id, currTimestamp], (err, results) => {
        if (err) {
            return next(err);
        }
        let newChannelID = results.insertId;
        if (req.body.private === true) {
            db.query(SQL_INSERT_INTO_CHANNEL_USER, [user.id, newChannelID], (err, results) => {
                if (err) {
                    return next(err);
                }
            });
        }
        db.query(SQL_SELECT_CHANNEL_BY_ID, [newChannelID], (err, rows) => {
            if (err) {
                return next(err);
            }
            let members = [];
            if (rows[0].private === 1) {
                members.push(user.id)
            }
            let channel = new Channel(rows[0].channelID, rows[0].name, rows[0].description,
                            rows[0].private, members, rows[0].createdAt, user, rows[0].editedAt);
            res.status(201);
            res.json(channel);
        });
    });
});

// Handle Endpoint: /v1/channels/{channelID}
app.get("/v1/channels/:channelID", (req, res, next) => {
    let user = checkUserAuth(req);
    if (user === false) {
        return res.status(401).send("Please sign in");
    }
    let valid = verifyUserInChannel(user.id, req.params.channelID, next);
    if (valid === false) {
        return res.status(403).send("Forbidden request. Not a part of this channel")
    }

    let messages = [];
    db.query(SQL_TOP_100_MESSAGES, [req.params.channelID], (err, rows) => {
       if (err) {
           return next(err);
       }
       rows.forEach((row) => {
           let message = new Message(row.id, row.channelID, row.body, row.createdAt, row.creator, row.editedAt);
           messages.push(message);
       });
    });
    res.json(messages);
});

app.post("/v1/channels/:channelID", (req, res, next) => {
    let user = checkUserAuth(req);
    if (user === false) {
        return res.status(401).send("Please sign in");
    }
    let valid = verifyUserInChannel(user.id, req.params.channelID, next);
    if (valid === false) {
        return res.status(403).send("Forbidden request. Not a part of this channel")
    }
    let dateNow = new Date().toISOString().slice(0, 19).replace('T', ' ');
    db.query(SQL_POST_MESSAGE, [req.params.channelID, req.body, dateNow, user.id, dateNow], (err, results) => {
        if (err) {
            return next(err);
        }
        let newMessageID = results.insertId;
        db.query(SQL_SELECT_MESSAGE_BY_ID, [newMessageID], (err, rows) => {
            if (err) {
                return next(err);
            }
            res.status(201);
            res.json(rows[0]);
        });
    });
});

app.patch("/v1/channels/:channelID", (req, res, next) => {
    let user = checkUserAuth(req);
    if (user === false) {
        return res.status(401).send("Please sign in");
    }
    let creatorStatus = checkUserIsCreator(user.id, next);
    if (creatorStatus === false) {
        return res.status(403).send("Error: You are not the creator of this channel");
    }
    db.query(SQL_UPDATE_CHANNEL_NAME_DESC, [req.body.name, req.body.description, req.params.channelID], (err, rows) => {
        if (err) {
            return next(err);
        }
    });
    db.query(SQL_SELECT_CHANNEL_BY_ID, [req.params.channelID], (err, rows) => {
        if (err) {
            return next(err);
        }
        res.json(rows[0]);
    });
});

app.delete("/v1/channels/:channelID", (req, res, next) => {
    let user = checkUserAuth(req);
    if (user === false) {
        return res.status(401).send("Please sign in");
    }
    let creatorStatus = checkUserIsCreator(user.id, next);
    if (creatorStatus === false) {
        return res.status(403).send("Error: You are not the creator of this channel");
    }
    db.query(SQL_DELETE_CHANNEL_AND_MESSAGES, [req.params.channelID], (err, rows) => {
        if (err) {
            return next(err);
        }
    });
    res.set("Content-Type", "text/plain");
    res.send("Successfully deleted channel and messages in that channel");
});

// Handle Endpoint: /v1/channels/{channelID}/members
app.post("/v1/channels/:channelID/members", (req, res, next) => {
    let user = checkUserAuth(req);
    if (user === false) {
        return res.status(401).send("Please sign in");
    }
    let creatorStatus = checkUserIsCreator(user.id, next);
    if (creatorStatus === false) {
        return res.status(403).send("Error: You are not the creator of this channel");
    }
    db.query(SQL_INSERT_INTO_CHANNEL_USER, [req.body.id, req.params.channelID], (err, rows) => {
        if (err) {
            return next(err);
        }
    });
    res.status(201);
    res.set("Content-Type", "text/plain");
    res.send("Successfully added user to channel");
});

app.delete("/v1/channels/:channelID/members", (req, res, next) => {
    let user = checkUserAuth(req);
    if (user === false) {
        return res.status(401).send("Please sign in");
    }
    let creatorStatus = checkUserIsCreator(user.id, next);
    if (creatorStatus === false) {
        return res.status(403).send("Error: You are not the creator of this channel");
    }
    db.query(SQL_DELETE_USER_FROM_CHANNEL, [req.params.channelID, user.id], (err, rows) => {
        if (err) {
            return next(err);
        }
    });
    res.status(200);
    res.set("Content-Type", "text/plain");
    res.send("Successfully removed user from channel");
});

// Handle Endpoint: /v1/messages/{messageID}
app.patch("/v1/messages/:messageID", (req, res, next) => {
    let user = checkUserAuth(req);
    if (user === false) {
        return res.status(401).send("Please sign in");
    }
    let creator = checkUserIsMessageCreator(user.id, req.params.messageID, next);
    if (creator === false) {
        return res.status(403).send("Error: You are not the creator of this channel");
    }
    db.query(SQL_UPDATE_MESSAGE, [req.body.body, new Date().toISOString().slice(0, 19).replace('T', ' '), req.params.messageID], (err, results) => {
        if (err) {
            return next(err);
        }
    });
    db.query(SQL_SELECT_MESSAGE_BY_ID, [], (err, rows) => {
        if (err) {
            return next(err);
        }
        res.status(200);
        res.json(rows[0])
    });
});

app.delete("/v1/messages/:messageID", (req, res, next) => {
    let user = checkUserAuth(req);
    if (user === false) {
        return res.status(401).send("Please sign in");
    }
    let creator = checkUserIsMessageCreator(user.id, req.params.messageID, next);
    if (creator === false) {
        return res.status(403).send("Error: You are not the creator of this channel");
    }
    db.query(SQL_DELETE_MESSAGE_BY_ID, [req.params.messageID], (err, results) => {
        if (err) {
            return next(err);
        }
    });
    res.set("Content-Type", "text/plain");
    res.send("Successfully deleted message");
});

app.use((err, req, res, next) => {
    if (err.stack) {
        console.error(err.stack);
    }
    res.status(500).send("Something bad happened. Sorry.");
});

app.listen(port, host, () => {
    console.log('server is listening at http://' + addr + '...');
});