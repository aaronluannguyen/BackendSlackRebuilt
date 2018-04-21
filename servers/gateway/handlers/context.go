package handlers

import (
	"github.com/challenges-aaronluannguyen/servers/gateway/sessions"
	"github.com/challenges-aaronluannguyen/servers/gateway/models/users"
)

//TODO: define a handler context struct that
//will be a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the key used for signing
//and verifying SessionIDs, the session store
//and the user store
type Context struct {
	signingKey		string
	sessionStore 	sessions.Store
	usersStore 		users.Store
}