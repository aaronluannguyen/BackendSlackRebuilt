class channel {
    constructor(
        id,
        name,
        description,
        private,
        members,
        createdAt,
        creator,
        editedAt
    ) {
        this.id = id;
        this.name = name;
        this.description = description;
        this.private = private;
        this.members = members;
        this.createdAt = createdAt;
        this.creator = creator;
        this.editedAt = editedAt;
    }
}

class message {
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
    }
}

module.exports = {
    channel,
    message
}