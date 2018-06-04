"use strict";

// Database Attribute ants
// users table
const user_id = "id";

// channels table
const channel_id = "channelID";
const channel_name = "channelName";
const channel_desc = "channelDescription";
const channel_private = "channelPrivate";
const channel_createAt = "channelCreatedAt";
const channel_creatorID = "channelCreatorUserID";
const channel_editedAt = "channelEditedAt";

// channel_user table
const cu_userID = "cuUserID";
const cu_channelID = "cuChannelID";

// messages table
const message_id = "mMessageID";
const message_chanID = "mChannelID";
const message_body = "mBody";
const message_createdAt = "mCreatedAt";
const message_creatorID = "mCreatorUserID";
const message_editedAt = "mEditedAt";

// message reaction table
const mrMessageID = "mrMessageID";
const mrUserID = "mrUserID";
const mrReactionCode = "mrReactionCode";

module.exports = {
    // SQL Commands
    SQL_SELECT_CHANNEL_BY_ID :                   "select * from channels where " + channel_id + "=?",
    SQL_UPDATE_CHANNEL_NAME_DESC :               "update channels set " + channel_name + "=?, " + channel_desc + "=?, " + channel_editedAt + "=? where " + channel_id + "=?",

    SQL_DELETE_CHANNEL :                         "delete channels" +
                                                 " from channels" +
                                                 " where " + channel_id + "=?",

    SQL_DELETE_CU :                              "delete channel_user" +
                                                 " from channel_user" +
                                                 " where " + cu_channelID + "=?",

    SQL_DELETE_CHANNEL_MESSAGES :                "delete messages" +
                                                 " from messages" +
                                                 " where " + message_chanID + "=?",

    SQL_SELECT_ALL_CHANNELS_FOR_USER :           "select * from channels c" +
                                                 " join channel_user cu on cu." + cu_channelID + " = c." + channel_id +
                                                 " join users u on u." + user_id  + "= cu." + cu_userID +
                                                 " where c." + channel_id + " in (" +
                                                        "select c." + channel_id + " from channels c" +
                                                        " join channel_user cu on cu." + cu_channelID + " = c." + channel_id +
                                                        " join users u on u." + user_id  + "= cu." + cu_userID +
                                                        " where c." + channel_private + "=? or cu." + cu_userID + "=?" +
                                                 ") order by c." + channel_id,

    SQL_SELECT_CHANNEL_MEMBERS :             "select * from channels c" +
                                                 " join channel_user cu on cu." + cu_channelID + " = c." + channel_id +
                                                 " join users u on u." + user_id  + "= cu." + cu_userID +
                                                 " where c." + channel_id + "=?" +
                                                 " order by u." + user_id,

    SQL_INSERT_NEW_CHANNEL :                     "insert into channels (" + channel_name + ", " + channel_desc + ", " + channel_private +
                                                 ", " + channel_createAt + ", " + channel_creatorID + ", " + channel_editedAt + ") " +
                                                 "values (?,?,?,?,?,?)",

    SQL_INSERT_INTO_CHANNEL_USER_BASE :          "insert into channel_user (" + cu_userID + ", " + cu_channelID + ") ",

    SQL_INSERT_INTO_CHANNEL_USER :               "insert into channel_user (" + cu_userID + ", " + cu_channelID + ") values(?,?)",

    SQL_DELETE_USER_FROM_CHANNEL :               "delete from channel_user where " + cu_channelID + "=? and " + cu_userID + "=?",

    SQL_SELECT_SPECIFIC_CHANNEL :                "select * from channels c" +
                                                 " join channel_user cu on cu." + cu_channelID + " = c." + channel_id +
                                                 " where c." + channel_id + "=?",


    SQL_TOP_100_MESSAGES:                        "select *, u1.id as id, u1.username as username, u1.firstName as firstName, u1.lastName as lastName, u1.photoURL as photoURL, u2.username as MRusername" +
                                                 " from (" +
                                                        "select * from messages m" +
                                                        " where m." + message_chanID + "=?" +
                                                        " order by m." + message_createdAt + " desc" +
                                                        " limit 100" +
                                                 ") as msg, users u1, users u2, message_reaction mr" +
                                                 " where msg." + message_id + " = mr." + mrMessageID +
                                                 " and u1." + user_id + " = msg." + message_creatorID +
                                                 " and u2." + user_id + " = mr." + mrUserID +
                                                 " order by msg." + message_createdAt,

    SQL_POST_MESSAGE :                           "insert into messages (" + message_chanID + ", " + message_body + ", " + message_createdAt +
                                                 ", " + message_creatorID + ", " + message_editedAt + ") values (?,?,?,?,?)",

    SQL_SELECT_MESSAGE_BY_ID :                   "select * from messages m" +
                                                 " join users u on u." + user_id  + " = m." + message_creatorID +
                                                 " where " + message_id + "=?",

    SQL_UPDATE_MESSAGE:                          "update messages set " + message_body + "=?, " + message_editedAt + "=? where " + message_id + "=?",

    SQL_DELETE_MESSAGE_BY_ID :                   "delete from messages where " + message_id + "=?",

    SQL_INSERT_INTO_MESSAGE_REACTION :           "insert into message_reaction (" + mrMessageID + ", " + mrUserID + ", " +  mrReactionCode + ") values (?,?,?)",

    SQL_GET_MESSAGE_WITH_REACTIONS :             "select * from message_reaction mr" +
                                                 " join users u on u." + user_id + "= mr." + mrUserID +
                                                 " where mr." + mrMessageID + "=?",

    SQL_POST_STAR_MESSAGE :                      "insert into star_message (smUserID, smMessageID) values(?,?)",

    SQL_GET_STAR_MESSAGES :                      "select * from star_message sm" +
                                                 " join users u on u." + user_id + " = sm.smUserID" +
                                                 " where sm.smUserID = ?",

    SQL_DELETE_STAR_MESSAGE :                    "delete from star_message" +
                                                 " where smUserID=? and smMessageID=?",
};