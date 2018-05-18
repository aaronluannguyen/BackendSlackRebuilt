"use strict";


let Channel = require("./models/channel");
let Message = require("./models/message");
let Constant = require("./models/sqlConstants");
const duplicateError = "ER_DUP_ENTRY";
const nonexistError = "ER_NO_REFERENCED";

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
app.get("/v1/channels", async (req, res, next) => {
    try {
        let user = checkUserAuth(req, res);
        if (!user) { return }
        let channels = await getChannelsForUser(db, false, user.id);
        res.json(channels);
    } catch (err) {
        next(err);
    }
});

app.post("/v1/channels", async (req, res, next) => {
    try {
        let user = checkUserAuth(req, res);
        if (!user) { return }
        if (!req.body.name) {
            res.set("Content-Type", "text/plain");
            return res.status(400).send("Please provide name for channel");
        }
        let timestamp = new Date().toISOString().slice(0, 19).replace('T', ' ');
        let result = await insertNewChannel(db, req.body.name, req.body.description,
                                            req.body.private, timestamp, user.id, req.body.members);
        if (result === duplicateError) {
            res.set("Content-Type", "text/plain");
            return res.status(400).send("Bad request: channel name already taken");
        }
        if (!result) {
            res.set("Content-Type", "text/plain");
            return res.status(500).send("server error: adding new channel");
        }
        let inserted = await newChannelInsertMembers(db, result.sqlCmd);
        if (!inserted) {
            let deleteStatus = await deleteChannelAndMessages(db, result.newChannelID);
            res.set("Content-Type", "text/plain");
            if (!deleteStatus) {
                return res.status(500).send("server error: deleting invalid channel add request");
            }
            return res.status(400).send("Bad Request: Tried inserting non-existent users");
        }
        let channel = await queryChannelMembers(db, result.newChannelID);
        if (!channel) {
            res.set("Content-Type", "text/plain");
            return res.status(500).send("server error: retrieving channel members");
        }
        res.status(201);
        res.json(channel);
    } catch (err) {
        next(err);
    }
});

// Handle Endpoint: /v1/channels/{channelID}
app.get("/v1/channels/:channelID", async (req, res, next) => {
    try {
        let user = checkUserAuth(req, res);
        if (!user) { return }
        let valid = await verifyUserInChannel(db, user.id, req.params.channelID);
        if (!valid) {
            res.set("Content-Type", "text/plain");
            return res.status(403).send("Forbidden request. Not a part of this channel");
        }
        let msgs = await queryTop100Msgs(db, req.params.channelID);
        res.json(msgs);
    } catch (err) {
        next(err);
    }
});

app.post("/v1/channels/:channelID", async (req, res, next) => {
    try {
        if (!validJSON) { return }
        let user = checkUserAuth(req, res);
        if (!user) { return }
        let valid = await verifyUserInChannel(db, user.id, req.params.channelID);
        if (!valid) {
            res.set("Content-Type", "text/plain");
            return res.status(403).send("Forbidden request. Not a part of this channel")
        }
        let dateNow = new Date().toISOString().slice(0, 19).replace('T', ' ');
        let newMessageID = await queryPostMessage(db, req.params.channelID, req.body.body, dateNow, user.id);
        let msg = await queryMessageByID(db, newMessageID);
        if (!msg) {
            res.set("Content-Type", "text/plain");
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
        let user = checkUserAuth(req, res);
        if (!user) { return }
        let creatorStatus = await checkUserIsCreator(db, user.id, req.params.channelID, res);
        if (!creatorStatus) { return }
        let update = await getOGChannelNameAndDesc(db, req.params.channelID);
        if (!update) {
            res.set("Content-Type", "text/plain");
            return res.status(400).send("No such channel");
        }
        if (checkIfNullEmpty(req.body.name) && checkIfNullEmpty(req.body.description)) {
            res.set("Content-Type", "text/plain");
            return res.status(400).send("Both Name and Description cannot be null or empty");
        }
        if (update.name === req.body.name && update.desc === req.body.description) {
            res.set("Content-Type", "text/plain");
            return res.status(403).send("No update necessary");
        }
        let newName = update.name;
        if (!checkIfNullEmpty(req.body.name)) {
            newName = req.body.name;
        }
        let newDesc = update.desc;
        if (!checkIfNullEmpty(req.body.description)) {
            newDesc = req.body.description;
        }
        let timestamp = new Date().toISOString().slice(0, 19).replace('T', ' ')
        let updated = await updateChannelNameAndDesc(db, newName, newDesc, req.params.channelID, timestamp);
        if (!updated) {
            res.set("Content-Type", "text/plain");
            return res.status(400).send("Error: Channel does not exist");
        }
        let channel = await queryChannelMembers(db, req.params.channelID);
        if (!channel) {
            res.set("Content-Type", "text/plain");
            return res.status(500).send("Error retrieving channel members");
        }
        res.status(201);
        res.json(channel);
    } catch (err) {
        next(err);
    }
});

app.delete("/v1/channels/:channelID", async (req, res, next) => {
    try {
        let user = checkUserAuth(req, res);
        if (!user) { return }
        let creatorStatus = await checkUserIsCreator(db, user.id, req.params.channelID, next, res);
        if (!creatorStatus) { return }
        let deleted = await deleteChannelAndMessages(db, req.params.channelID);
        res.set("Content-Type", "text/plain");
        if (deleted === false) {
            res.status(400).send("Bad request: message does not exist");
        }
        res.send("Successfully deleted channel and messages in that channel");
    } catch (err) {
        next(err);
    }
});

// Handle Endpoint: /v1/channels/{channelID}/members
app.post("/v1/channels/:channelID/members", async (req, res, next) => {
    try {
        let user = checkUserAuth(req, res);
        if (!user) { return }
        let creatorStatus = await checkUserIsCreator(db, user.id, req.params.channelID, res);
        if (!creatorStatus) { return }
        let added = await queryAddUserToChannel(db, req.body.id, req.params.channelID);
        if (added === duplicateError) {
            res.set("Content-Type", "text/plain");
            return res.status(400).send("Bad request: member already added to channel");
        }
        if (added === nonexistError) {
            res.set("Content-Type", "text/plain");
            return res.status(400).send("Bad request: member does not exist");
        }
        if (!added) {
            res.set("Content-Type", "text/plain");
            return res.status(500).send("Server error adding user to channel");
        }
        res.status(201);
        res.set("Content-Type", "text/plain");
        res.send("Successfully added user to channel");
    } catch (err) {
        next(err);
    }
});

app.delete("/v1/channels/:channelID/members", async (req, res, next) => {
    try {
        let user = checkUserAuth(req, res);
        if (!user) { return }
        let creatorStatus = await checkUserIsCreator(db, user.id, req.params.channelID, res);
        if (!creatorStatus) { return }
        let deleted = await queryDeleteUserFromChannel(db, req.params.channelID, req.body.id);
        if (!deleted) {
            res.set("Content-Type", "text/plain");
            res.status(400).send("User does not exist in channel");
        }
        res.set("Content-Type", "text/plain");
        res.status(200).send("Successfully removed user from channel");
    } catch (err) {
        next(err);
    }
});

// Handle Endpoint: /v1/messages/{messageID}
app.patch("/v1/messages/:messageID", async (req, res, next) => {
    try {
        let user = checkUserAuth(req, res);
        if (!user) { return }
        let creator = await checkUserIsMessageCreator(db, user.id, req.params.messageID, res);
        if (!creator) { return }
        let updated = await queryUpdateMsg(db, req.body.body,
                                new Date().toISOString().slice(0, 19).replace('T', ' '), req.params.messageID);
        if (!updated) {
            res.set("Content-Type", "text/plain");
            return res.status(500).send("Server error updating message");
        }
        let msg = await queryMessageByID(db, req.params.messageID);
        if (!msg) {
            res.set("Content-Type", "text/plain");
            return res.status(400).send("Message does not exist");
        }
        res.json(msg);
    } catch (err) {
        next(err);
    }
});

app.delete("/v1/messages/:messageID", async (req, res, next) => {
    try {
        let user = checkUserAuth(req, res);
        if (!user) { return }
        let creator = await checkUserIsMessageCreator(db, user.id, req.params.messageID, res);
        if (!creator) { return }
        let deleted = await queryDeleteMessage(db, req.params.messageID);
        if (!deleted) {
            res.set("Content-Type", "text/plain");
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
    res.set("Content-Type", "text/plain");
    res.status(500).send("Server Error...");
});

app.listen(port, host, () => {
    console.log('server is listening at http://' + addr + '...');
});

//checkUserAuth checks the x-user header to make sure the user is signed in
//if user is signed in, the user, as json, is returned
//if not signed in, returns false
function checkUserAuth(req, res) {
    let userJson = req.get("X-User");
    if (!userJson) {
        res.set("Content-Type", "text/plain");
        res.status(401).send("Error: Please sign in");
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
            if (!rows[0].channelPrivate) {
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
function checkUserIsCreator(db, userID, channelID, res) {
    return new Promise((resolve, reject) => {
        db.query(Constant.SQL_SELECT_CHANNEL_BY_ID, [channelID], (err, rows) => {
            if (err) {
                reject(err);
            }
            if (!rows || rows.length === 0) {
                res.set("Content-Type", "text/plain");
                res.status(403).send("Error: You are not the creator of this channel");
                return resolve(false);
            }
            if (!rows[0].channelPrivate) {
                res.set("Content-Type", "text/plain");
                res.status(400).send("Bad request: This is a public channel");
                return resolve(false)
            }
            if (rows[0].channelCreatorUserID === userID) {
                return resolve(true);
            }
            res.set("Content-Type", "text/plain");
            res.status(403).send("Error: You are not the creator of this channel");
            resolve(false);
        });
    });
}

//checkUserIsMessageCreator checks if the user is the creator of a specified message
//and returns true or false accordingly
function checkUserIsMessageCreator(db, userID, messageID, res) {
    return new Promise((resolve, reject) => {
        db.query(Constant.SQL_SELECT_MESSAGE_BY_ID, [messageID], (err, rows) => {
            if (err) {
                reject(err);
            }
            if (rows.length === 0) {
                res.set("Content-Type", "text/plain");
                res.status(403).send("Error: You are not the creator of this message");
                return resolve(false);
            }
            if (rows[0].mCreatorUserID === userID) {
                return resolve(true);
            }
            res.set("Content-Type", "text/plain");
            res.status(403).send("Error: You are not the creator of this message");
            return resolve(false);
        });
    });
}

// getChannelsForUser returns all channels available to user
function getChannelsForUser(db, bool, userID) {
    return new Promise((resolve, reject) => {
        let channels = [];
        db.query(Constant.SQL_SELECT_ALL_CHANNELS_FOR_USER, [bool, userID], (err, rows) => {
            if (err) {
                reject(err);
            }
            if (rows.length === 0) {
                return resolve(channels);
            }
            let members = [];
            let currChannelID = rows[0].channelID;
            let channel = new Channel(rows[0].channelID, rows[0].channelName, rows[0].channelDescription,
                rows[0].channelPrivate, members, rows[0].channelCreatedAt,
                {}, rows[0].channelEditedAt);
            rows.forEach((row) => {
                if (row.channelID !== currChannelID) {
                    if (!channel.private) {
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
            if (channel.private === 0) {
                channel.members = [];
            }
            channels.push(channel);
            return resolve(channels);
        })
    });
}

//insertNewChannel inserts a new channel and returns a promise containing the sql command for adding members
function insertNewChannel(db, name, descr, isPrivate, date, userID, members) {
    return new Promise((resolve, reject) => {
        db.query(Constant.SQL_INSERT_NEW_CHANNEL, [name, descr, isPrivate, date, userID, date], (err, results) => {
            if (err) {
                if (err.message.startsWith(duplicateError)) {
                    return resolve(duplicateError);
                }
                reject(err);
            }
            if (results.affectedRows === 0) {
                return resolve(false);
            }
            let newChannelID = results.insertId;
            let allMembersInsert = "values (" + userID + ", " + newChannelID + "), ";
            let membersToAdd = [];
            if (members) {
                membersToAdd = members;
            }
            membersToAdd.forEach((member) => {
                let memberObj = "(" + member.id + ", " + newChannelID + "), ";
                allMembersInsert += memberObj;
            });
            let trimmedAllMembersInsert = allMembersInsert.slice(0, allMembersInsert.length - 2);
            trimmedAllMembersInsert += ";";
            let membersInsertSQL = Constant.SQL_INSERT_INTO_CHANNEL_USER_BASE + trimmedAllMembersInsert;
            return resolve({newChannelID: newChannelID, sqlCmd: membersInsertSQL});
        });
    });
}

//newChannelInsertMembers inserts given users at channel creation and
// returns a promise indicating status of insertion
function newChannelInsertMembers(db, sql) {
    return new Promise((resolve, reject) => {
        db.query(sql, (err, rows) => {
            if (err) {
                return resolve(false);
                reject(err);
            }
            if (rows.affectedRows === 0) {
                return resolve(false)
            }
            return resolve(true);
        });
    });
}

//queryDeleteUserFromChannel returns a promise whether the user was successfully deleted from a channel
function queryDeleteUserFromChannel(db, channelID, userID) {
    return new Promise((resolve, reject) => {
        db.query(Constant.SQL_DELETE_USER_FROM_CHANNEL, [channelID, userID], (err, results) => {
            if (err) {
                reject(err);
            }
            if (results.affectedRows === 0) {
                return resolve(false);
            }
            return resolve(true);
        });
    });
}

//queryAddUserToChannel returns a promise that indicates if a user was successfully added to a channel
function queryAddUserToChannel(db, userID, channelID) {
    return new Promise((resolve, reject) => {
        db.query(Constant.SQL_INSERT_INTO_CHANNEL_USER, [userID, channelID], (err, results) => {
            if (err) {
                if (err.message.startsWith(duplicateError)) {
                    return resolve(duplicateError);
                }
                if (err.message.startsWith(nonexistError)) {
                    return resolve(nonexistError);
                }
                reject(err);
            }
            if (results.affectedRows === 0) {
                return resolve(false);
            }
            resolve(true);
        });
    });
}

//getOGChannelNameAndDesc returns a promise containing the current Channel Name and Description
function getOGChannelNameAndDesc(db, channelID) {
    return new Promise((resolve, reject) => {
        db.query(Constant.SQL_SELECT_CHANNEL_BY_ID, [channelID], (err, rows) => {
            if (err) {
                reject(err);
            }
            if (!rows || rows.length === 0) {
                return resolve(false);
            }
            resolve({name: rows[0].channelName, desc: rows[0].channelDescription});
        });
    });
}

//checkIfNullEmpty checks if channel name or desc is empty/null
function checkIfNullEmpty(obj) {
    return (!obj || obj === "");
}

//updateChannelNameAndDesc returns a promise if successfully updated channel
function updateChannelNameAndDesc(db, name, desc, channelID, date) {
    return new Promise((resolve, reject) => {
        db.query(Constant.SQL_UPDATE_CHANNEL_NAME_DESC, [name, desc, channelID, date], (err, results) => {
            if (err) {
                reject(err);
            }
            if (results.affectedRows === 0) {
                return resolve(false);
            }
            resolve(true);
        });
    });
}

//queryChannelMembers returns a promise containing rows for finding members in a channel
function queryChannelMembers(db, channelID) {
    return new Promise((resolve, reject) => {
        db.query(Constant.SQL_SELECT_CHANNEL_MEMBERS, [channelID], (err, rows) => {
            if (err) {
                reject(err);
            }
            if (!rows || rows.length === 0) {
                return resolve(false)
            }
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
            if (!rows[0].channelPrivate) {
                members = [];
            }
            let channel = new Channel(rows[0].channelID, rows[0].channelName, rows[0].channelDescription,
                rows[0].channelPrivate, members, rows[0].channelCreatedAt, creator, rows[0].channelEditedAt);
            return resolve(channel);
        });
    });
}

//queryPostMessage posts a message and returns a promise containing the message ID
function queryPostMessage(db, channelID, body, date, userID) {
    return new Promise((resolve, reject) => {
        db.query(Constant.SQL_POST_MESSAGE, [channelID, body, date, userID, date], (err, results) => {
            if (err) {
                reject(err);
            }
            if (results.affectedRows === 0) {
                return resolve(false);
            }
            resolve(results.insertId);
        });
    });
}

//queryMessageByID returns a promise containing a message model
function queryMessageByID(db, msgID) {
    return new Promise((resolve, reject) => {
        db.query(Constant.SQL_SELECT_MESSAGE_BY_ID,[msgID], (err, rows) => {
            if (err) {
                reject(err);
            }
            if (!rows || rows.length === 0) {
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
function queryTop100Msgs(db, channelID) {
    return new Promise((resolve, reject) => {
        let messages = [];
        db.query(Constant.SQL_TOP_100_MESSAGES, [channelID], (err, rows) => {
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
function queryUpdateMsg(db, body, date, messageID) {
    return new Promise((resolve, reject) => {
        db.query(Constant.SQL_UPDATE_MESSAGE, [body, date, messageID], (err, results) => {
            if (err) {
                reject(err);
            }
            if (!results) {
                return resolve(false);
            }
            resolve(true);
        });
    });
}

//queryDeleteMessage queries to delete a message and returns a promise indicating status of deletion
function queryDeleteMessage(db, messageID) {
    return new Promise((resolve, reject) => {
        db.query(Constant.SQL_DELETE_MESSAGE_BY_ID, [messageID], (err, results) => {
            if (err) {
                reject(err);
            }
            if (!results) {
                return resolve(false);
            }
            resolve(true);
        });
    });
}

//deleteChannelAndMessages queries to delete all data related to given channelID such as
//from the channels, channel_user, and messages tables from the data base
function deleteChannelAndMessages(db, channelID) {
    return new Promise((resolve, reject) => {
        resolve(true);
        db.query(Constant.SQL_DELETE_CHANNEL_MESSAGES, [channelID], (err, results) => {
            if (err) {
                reject(err);
            }
            if (!results) {
                return resolve(false);
            }
            db.query(Constant.SQL_DELETE_CU, [channelID], (err, results) => {
                if (err) {
                    reject(err);
                }
                if (!results) {
                    return resolve(false);
                }
                db.query(Constant.SQL_DELETE_CHANNEL, [channelID], (err, results) => {
                    if (err) {
                        reject(err);
                    }
                    if (!results) {
                        return resolve(false);
                    }
                });
            });
        });
    });
}