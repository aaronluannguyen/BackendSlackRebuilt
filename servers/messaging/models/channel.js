"use strict";

module.exports = class channel {
    constructor(
        id,
        name,
        description,
        isPrivate,
        members,
        createdAt,
        creator,
        editedAt
    ) {
        this.id = id;
        this.name = name;
        this.description = description;
        this.isPrivate = isPrivate;
        this.members = members;
        this.createdAt = createdAt;
        this.creator = creator;
        this.editedAt = editedAt;
    }
};