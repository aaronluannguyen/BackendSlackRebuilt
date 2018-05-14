export class channel {
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

    toJSON() {
        return {
            id: this.id,
            name: this.name,
            description: this.description,
            private: this.private,
            members: this.members,
            createdAt: this.createdAt,
            creator: this.creator,
            editedAt: this.editedAt
        };
    }
}

export class message {
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

    toJSON() {
        return {
            id: this.id,
            channelID: this.channelID,
            body: this.body,
            createdAt: this.createdAt,
            creator: this.creator,
            editedAt: this.editedAt
        };
    }
}