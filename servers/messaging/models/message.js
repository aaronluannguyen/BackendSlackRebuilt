"use strict";

module.exports = class message {
    constructor (
        id,
        channelID,
        body,
        createdAt,
        creator,
        editedAt
    ) {
        this.id = id;
        this.channelID = channelID;
        this.body = body;
        this.createdAt = createdAt;
        this.creator = creator;
        this.editedAt = editedAt;
        this.reactions = [];
    }
};