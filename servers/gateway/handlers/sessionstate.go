package handlers

import (
	"github.com/challenges-aaronluannguyen/servers/gateway/models/users"
	"time"
)

//TODO: define a session state struct for this web server
//see the assignment description for the fields you should include
//remember that other packages can only see exported fields!
type SessionState struct {
	time 	time.Time 		`json: "time"`
	user 	*users.User 		`json: "user"`
}