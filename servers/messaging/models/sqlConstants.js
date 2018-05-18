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

module.exports = {
    // SQL Commands
    SQL_SELECT_CHANNEL_BY_ID :                   "select * from channels where " + channel_id + "=?",
    SQL_UPDATE_CHANNEL_NAME_DESC :               "update channels set " + channel_name + "=?, " + channel_desc + "=? where " + channel_id + "=?",

    SQL_ALTER_TABLE_BEFORE_CHANNEL_DELETE :      "alter table channels add constraint channel_user_ibfk_2" +
                                                 " foreign key (" + channel_id +") references channel_user (" + cu_channelID + ");",

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

    SQL_TOP_100_MESSAGES :                       "select * from channels c" +
                                                 " join messages m on m. " + message_chanID + " = c." + channel_id +
                                                 " join users u on u." + user_id  + " = m." + message_creatorID +
                                                 " where c." + channel_id + " = ?" +
                                                 " order by m." + message_createdAt + " desc" +
                                                 " limit 100",

    SQL_POST_MESSAGE :                           "insert into messages (" + message_chanID + ", " + message_body + ", " + message_createdAt +
                                                 ", " + message_creatorID + ", " + message_editedAt + ") values (?,?,?,?,?)",

    SQL_SELECT_MESSAGE_BY_ID :                   "select * from messages m" +
                                                 " join users u on u." + user_id  + " = m." + message_creatorID +
                                                 " where " + message_id + "=?",

    SQL_UPDATE_MESSAGE:                          "update messages set " + message_body + "=?, " + message_editedAt + "=? where " + message_id + "=?",

    SQL_DELETE_MESSAGE_BY_ID :                   "delete from messages where " + message_id + "=?",
};