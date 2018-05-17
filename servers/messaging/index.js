"use strict";


let Channel = require("./models/channel");
let Message = require("./models/message");
let Constant = require("./models/sqlConstants");

const mysql = require("mysql");
const express = require("express");
const app = express();

const addr = process.env.ADDR || ":80";
const [host, port] = addr.split(":");

let db = mysql.createPool({
    host: process.env.MYSQL_ADDR,
    database: process.env.MYSQL_DATABASE,
    user: "root",
    password: process.env.MYSQL_ROOT_PASSWORD
});

app.use(express.json());

// Handle Endpoint: /v1/channels
app.get("/v1/channels", (req, res, next) => {
    let user = checkUserAuth(req);
    if (user === false) {
        return res.status(401).send("Please sign in");
    }

    let channels = [];
    db.query(Constant.SQL_SELECT_ALL_CHANNELS_FOR_USER, [false, user.id], (err, rows) => {
        if (err) {
            return next(err);
        }
        if (rows.length === 0) {
            return res.json(channels);
        }
        let members = [];
        let currChannelID = rows[0].channelID;
        let channel = new Channel(rows[0].channelID, rows[0].channelName, rows[0].channelDescription,
                                    rows[0].channelPrivate, members, rows[0].channelCreatedAt,
                                    {}, rows[0].channelEditedAt);

        rows.forEach((row) => {
            if (row.channelID !== currChannelID) {
                if (channel.isPrivate === 0) {
                    channel.members = [];
                }
                channels.push(channel);
                members = [{id: row.id, userName: row.username, firstName: row.firstName,
                            lastName: row.lastName, photoURL: row.photoURL}];
                channel = new Channel(row.channelID, row.channelName, row.channelDescription,
                    row.channelPrivate, members, row.channelCreatedAt,
                    {}, row.channelEditedAt);
                if (row.channelCreatorUserID === row.id) {
                    channel.creator = {id: row.channelCreatorUserID, userName: row.username,
                        firstName: row.firstName, lastName: row.lastName, photoURL: row.photoURL};
                }
                currChannelID = row.channelID;
            } else {
                let member = {id: row.id, userName: row.username, firstName: row.firstName,
                    lastName: row.lastName, photoURL: row.photoURL};
                members.push(member);
                if (row.channelCreatorUserID === row.id) {
                    channel.creator = {id: row.channelCreatorUserID, userName: row.username,
                        firstName: row.firstName, lastName: row.lastName, photoURL: row.photoURL};
                }
            }
        });
        if (channel.isPrivate === 0) {
            channel.members = [];
        }
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
    db.query(Constant.SQL_INSERT_NEW_CHANNEL, [req.body.name, req.body.description, req.body.private,
                currTimestamp, user.id, currTimestamp], (err, results) => {
        if (err) {
            return next(err);
        }
        let newChannelID = results.insertId;
        let allMembersInsert = "values (" + user.id + ", " + newChannelID + "), ";
        let membersToAdd = [];
        if (req.body.members) {
            membersToAdd = req.body.members;
        }
        membersToAdd.forEach((member) => {
            let memberObj = "(" + member.id + ", " + newChannelID + "), ";
            allMembersInsert += memberObj;
        });
        let trimmedAllMembersInsert = allMembersInsert.slice(0, allMembersInsert.length - 2);
        trimmedAllMembersInsert += ";";
        let membersInsertSQL = Constant.SQL_INSERT_INTO_CHANNEL_USER_BASE + trimmedAllMembersInsert;
        db.query(membersInsertSQL, (err, results) => {
            if (err) {
                return next(err);
            }
            db.query(Constant.SQL_SELECT_CHANNEL_MEMBERS, [newChannelID], (err, rows) => {
                if (err) {
                    return next(err);
                }
                let members = [];
                if (rows[0].channelPrivate === 1) {
                    rows.forEach((row) => {
                        let member = {id: row.id, userName: row.username, firstName: row.firstName,
                            lastName: row.lastName, photoURL: row.photoURL};
                        members.push(member);
                    });
                }
                let channel = new Channel(rows[0].channelID, rows[0].channelName, rows[0].channelDescription,
                    rows[0].channelPrivate, members, rows[0].channelCreatedAt, user, rows[0].channelEditedAt);
                res.status(201);
                res.json(channel);
            });
        });
    });
});

// Handle Endpoint: /v1/channels/{channelID}
app.get("/v1/channels/:channelID", async (req, res, next) => {
    try {
        let user = checkUserAuth(req);
        if (user === false) {
            return res.status(401).send("Please sign in");
        }
        let valid = await verifyUserInChannel(db, user.id, req.params.channelID);
        if (valid === false) {
            return res.status(403).send("Forbidden request. Not a part of this channel");
        }
        let msgs = await queryTop100Msgs(db, Constant.SQL_TOP_100_MESSAGES, req.params.channelID);
        res.json(msgs);
    } catch (err) {
        next(err);
    }
});

app.post("/v1/channels/:channelID", async (req, res, next) => {
    try {
        let user = checkUserAuth(req);
        if (user === false) {
            return res.status(401).send("Please sign in");
        }
        let valid = await verifyUserInChannel(db, user.id, req.params.channelID);
        if (valid === false) {
            return res.status(403).send("Forbidden request. Not a part of this channel")
        }
        let dateNow = new Date().toISOString().slice(0, 19).replace('T', ' ');
        let results = await queryPostMessage(db, Constant.SQL_POST_MESSAGE, req.params.channelID, req.body.body, dateNow, user.id);
        let newMessageID = results.insertId;
        let msg = await queryMessageByID(db, Constant.SQL_SELECT_MESSAGE_BY_ID, newMessageID);
        if (!msg) {
            return res.status(400).send("Message does not exist");
        }
        res.status(201);
        res.json(msg);
    } catch (err) {
        next(err);
    }
});

app.patch("/v1/channels/:channelID", async (req, res, next) => {
    try {
        let user = checkUserAuth(req);
        if (user === false) {
            return res.status(401).send("Please sign in");
        }
        let creatorStatus = await checkUserIsCreator(user.id, req.params.channelID);
        if (creatorStatus === false) {
            return res.status(403).send("Error: You are not the creator of this channel");
        }
        let update = await getOGChannelNameAndDesc(Constant.SQL_SELECT_CHANNEL_BY_ID, req.params.channelID);
        if (!update) {
            return res.status(400).send("No such channel");
        }
        let newName = update.name;
        if (!checkIfChannelUpdateNullEmpty(req.body.name)) {
            newName = req.body.name;
        }
        let newDesc = update.desc;
        if (!checkIfChannelUpdateNullEmpty(req.body.description)) {
            newDesc = req.body.description;
        }
        await updateChannelNameAndDesc(Constant.SQL_UPDATE_CHANNEL_NAME_DESC, newName,
            newDesc, req.params.channelID);
        let rows = await queryChannelMembers(Constant.SQL_SELECT_CHANNEL_MEMBERS, req.params.channelID);
        let members = [];
        let creator = {};
        rows.forEach((row) => {
            let member = {id: row.id, userName: row.username, firstName: row.firstName,
                lastName: row.lastName, photoURL: row.photoURL};
            members.push(member);
            if (row.channelCreatorUserID === row.id) {
                creator = {id: row.id, userName: row.username, firstName: row.firstName,
                    lastName: row.lastName, photoURL: row.photoURL};
            }
        });
        if (rows[0].channelPrivate === 1) {
            members = [];
        }
        let channel = new Channel(rows[0].channelID, rows[0].channelName, rows[0].channelDescription,
            rows[0].channelPrivate, members, rows[0].channelCreatedAt, creator, rows[0].channelEditedAt);
        res.status(201);
        res.json(channel);
    } catch (err) {
        next(err);
    }
});

app.delete("/v1/channels/:channelID", async (req, res, next) => {
    try {
        let user = checkUserAuth(req);
        if (user === false) {
            return res.status(401).send("Please sign in");
        }
        let creatorStatus = await checkUserIsCreator(user.id, req.params.channelID, next);
        if (creatorStatus === false) {
            return res.status(403).send("Error: You are not the creator of this channel");
        }
    } catch (err) {
        next(err);
    }

    db.query(Constant.SQL_ALTER_TABLE_BEFORE_CHANNEL_DELETE, (err, rows) => {
        if (err) {
            return next(err);
        }
    });
    db.query(Constant.SQL_DELETE_CHANNEL_AND_MESSAGES, [req.params.channelID], (err, rows) => {
        if (err) {
            return next(err);
        }
        res.set("Content-Type", "text/plain");
        res.send("Successfully deleted channel and messages in that channel");
    });
});

// Handle Endpoint: /v1/channels/{channelID}/members
app.post("/v1/channels/:channelID/members", async (req, res, next) => {
    try {
        let user = checkUserAuth(req);
        if (user === false) {
            return res.status(401).send("Please sign in");
        }
        let creatorStatus = await checkUserIsCreator(user.id, req.params.channelID, next);
        if (creatorStatus === false) {
            return res.status(403).send("Error: You are not the creator of this channel");
        }
        await queryAddUserToChannel(db, Constant.SQL_INSERT_INTO_CHANNEL_USER,
                                    req.body.id, req.params.channelID);
        res.status(201);
        res.set("Content-Type", "text/plain");
        res.send("Successfully added user to channel");
    } catch (err) {
        next(err);
    }
});

app.delete("/v1/channels/:channelID/members", async (req, res, next) => {
    try {
        let user = checkUserAuth(req);
        if (user === false) {
            return res.status(401).send("Please sign in");
        }
        let creatorStatus = await checkUserIsCreator(user.id, req.params.channelID);
        if (!creatorStatus) {
            return res.status(403).send("Error: You are not the creator of this channel");
        }
        await queryDeleteUserFromChannel(db, Constant.SQL_DELETE_USER_FROM_CHANNEL,
                                        req.params.channelID, req.body.id);
        res.status(200);
        res.set("Content-Type", "text/plain");
        res.send("Successfully removed user from channel");
    } catch (err) {
        next(err);
    }
});

// Handle Endpoint: /v1/messages/{messageID}
app.patch("/v1/messages/:messageID", async (req, res, next) => {
    try {
        let user = checkUserAuth(req);
        if (user === false) {
            return res.status(401).send("Please sign in");
        }
        let creator = await checkUserIsMessageCreator(db, user.id, req.params.messageID);
        if (creator === false) {
            return res.status(403).send("Error: You are not the creator of this message");
        }
        await queryUpdateMsg(db, Constant.SQL_UPDATE_MESSAGE, req.body.body,
                                new Date().toISOString().slice(0, 19).replace('T', ' '), req.params.messageID);
        let msg = await queryMessageByID(db, Constant.SQL_SELECT_MESSAGE_BY_ID, req.params.messageID);
        if (!msg) {
            return res.status(400).send("Message does not exist");
        }
        res.json(msg);
    } catch (err) {
        next(err);
    }
});

app.delete("/v1/messages/:messageID", async (req, res, next) => {
    try {
        let user = checkUserAuth(req);
        if (user === false) {
            return res.status(401).send("Please sign in");
        }
        let creator = await checkUserIsMessageCreator(db, user.id, req.params.messageID);
        if (creator === false) {
            return res.status(403).send("Error: You are not the creator of this message");
        }
        let deleted = await queryDeleteMessage(db, Constant.SQL_DELETE_MESSAGE_BY_ID, req.params.messageID);
        if (!deleted) {
            return res.status(400).send("Message does not exist")
        }
        res.set("Content-Type", "text/plain");
        res.send("Successfully deleted message");
    } catch (err) {
        next(err);
    }
});

app.use((err, req, res, next) => {
    if (err.stack) {
        console.error(err.stack);
    }
    res.status(500).send("Server Error...");
});

app.listen(port, host, () => {
    console.log('server is listening at http://' + addr + '...');
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
function verifyUserInChannel(db, userID, channelID) {
    return new Promise((resolve, reject) => {
        db.query(Constant.SQL_SELECT_SPECIFIC_CHANNEL, [channelID], (err, rows) => {
            if (err) {
                reject(err);
            }
            if (rows.length === 0) {
                return resolve(false);
            }
            if (rows[0].channelPrivate === false) {
                return resolve(true);
            }
            rows.forEach((row) => {
                if (row.cuUserID === userID) {
                    return resolve(true);
                }
            });
            return resolve(false);
        });
    });
}

//checkUserIsCreator checks if the user is the creator of a specified channel
//and returns true or false accordingly
function checkUserIsCreator(userID, channelID) {
    return new Promise((resolve, reject) => {
        db.query(Constant.SQL_SELECT_CHANNEL_BY_ID, [channelID], (err, rows) => {
            if (err) {
                reject(err);
            }
            if (rows.length === 0) {
                return resolve(false);
            }
            resolve(rows[0].channelCreatorUserID === userID);
        });
    });
}

//checkUserIsMessageCreator checks if the user is the creator of a specified message
//and returns true or false accordingly
function checkUserIsMessageCreator(db, userID, messageID) {
    return new Promise((resolve, reject) => {
        db.query(Constant.SQL_SELECT_MESSAGE_BY_ID, [messageID], (err, rows) => {
            if (err) {
                reject(err);
            }
            if (rows.length === 0) {
                return resolve(false)
            }
            resolve(rows[0].mCreatorUserID === userID);
        });
    });
}

//queryDeletedUserFromChannel returns a promise whether the user was successfully deleted from a channel
function queryDeleteUserFromChannel(db, sql, channelID, userID) {
    return new Promise((resolve, reject) => {
        db.query(sql, [channelID, userID], (err, rows) => {
            if (err) {
                reject(err);
            }
            resolve(true);
        });
    });
}

//queryAddUserToChannel returns a promise that indicates if a user was successfully added to a channel
function queryAddUserToChannel(db, sql, userID, channelID) {
    return new Promise((resolve, reject) => {
        db.query(sql, [userID, channelID], (err, rows) => {
            if (err) {
                reject(err);
            }
            resolve(true);
        });
    });
}

//getOGChannelNameAndDesc returns a promise containing the current Channel Name and Description
function getOGChannelNameAndDesc(sql, channelID) {
    return new Promise((resolve, reject) => {
        db.query(sql, [channelID], (err, rows) => {
            if (err) {
                reject(err);
            }
            if (rows.length === 0) {
                return resolve(false);
            }
            resolve({name: rows[0].channelName, desc: rows[0].channelDescription});
        });
    });
}

//checkIfChannelUpdateNullEmpty checks if channel name or desc is empty/null
function checkIfChannelUpdateNullEmpty(obj) {
    return (!obj || obj === "");
}

//updateChannelNameAndDesc returns a promise if successfully updated channel
function updateChannelNameAndDesc(sql, name, desc, channelID) {
    return new Promise((resolve, reject) => {
        db.query(sql, [name, desc, channelID], (err, rows) => {
            if (err) {
                reject(err);
            }
            resolve(true);
        });
    });
}

//queryChannelMembers returns a promise containing rows for finding members in a channel
function queryChannelMembers(sql, channelID) {
    return new Promise((resolve, reject) => {
        db.query(sql, [channelID], (err, rows) => {
            if (err) {
                reject(err);
            }
            resolve(rows);
        });
    });
}

//queryPostMessage posts a message and returns a promise containing the message ID
function queryPostMessage(db, sql, channelID, body, date, userID) {
    return new Promise((resolve, reject) => {
        db.query(sql, [channelID, body, date, userID, date], (err, results) => {
            if (err) {
                reject(err);
            }
            resolve(results);
        });
    });
}

//queryMessageByID returns a promise containing a message model
function queryMessageByID(db, sql, msgID) {
    return new Promise((resolve, reject) => {
        db.query(sql,[msgID], (err, rows) => {
            if (err) {
                reject(err);
            }
            if (rows.length === 0) {
                return resolve(false);
            }

            let creator = {id: rows[0].id, userName: rows[0].username, firstName: rows[0].firstName,
                lastName: rows[0].lastName, photoURL: rows[0].photoURL};
            let message = new Message(rows[0].mMessageID, rows[0].mChannelID, rows[0].mBody,
                rows[0].mCreatedAt, creator, rows[0].mEditedAt);
            resolve(message);
        });
    });
}

//queryTop100Msgs returns a promise containing the most recent 100 messages in a channel
function queryTop100Msgs(db, sql, channelID) {
    return new Promise((resolve, reject) => {
        let messages = [];
        db.query(sql, [channelID], (err, rows) => {
            if (err) {
                reject(err);
            }
            if (rows.length === 0) {
                return resolve(messages);
            }
            rows.forEach((row) => {
                let creator = {id: row.id, userName: row.username, firstName: row.firstName,
                    lastName: row.lastName, photoURL: row.photoURL};
                let message = new Message(row.mMessageID, row.mChannelID, row.mBody,
                    row.mCreatedAt, creator, row.mEditedAt);
                messages.push(message);
            });
            resolve(messages);
        });
    });
}

//queryUpdateMsg queries to update a message's body and returns a promise confirming that msg was updated
function queryUpdateMsg(db, sql, body, date, messageID) {
    return new Promise((resolve, reject) => {
        db.query(sql, [body, date, messageID], (err, results) => {
            if (err) {
                reject(err);
            }
            resolve(true);
        });
    });
}

//queryDeleteMessage queries to delete a message and returns a promise indicating status of deletion
function queryDeleteMessage(db, sql, messageID) {
    return new Promise((resolve, reject) => {
        db.query(sql, [messageID], (err, rows) => {
            if (err) {
                reject(err);
            }
            if (rows.length === 0) {
                return resolve(false);
            }
            resolve(true);
        });
    });
}