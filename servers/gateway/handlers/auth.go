package handlers

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/challenges-aaronluannguyen/servers/gateway/models/users"
	"github.com/challenges-aaronluannguyen/servers/gateway/sessions"
)

//TODO: define HTTP handler functions as described in the
//assignment description. Remember to use your handler context
//struct as the receiver on these functions so that you have
//access to things like the session store and user store.
func (ctx *Context) UsersHandler(w http.ResponseWriter, r *http.Request) {

}

func (ctx *Context) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (ctx *Context) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		contentType := r.Header.Get("Content-type")
		if contentType != contentTypeJSON {
			http.Error(w, fmt.Sprintf("request body must be in JSON"), http.StatusUnsupportedMediaType)
			return
		}
		cred := users.Credentials{}
		if err := json.NewDecoder(r.Body).Decode(cred); err != nil {
			http.Error(w, fmt.Sprintf("error decoding into credentials: %v", err), http.StatusBadRequest)
			return
		}
		user, err := ctx.usersStore.GetByEmail(cred.Email)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid credentials"), http.StatusUnauthorized)
			return
		}
		if err := user.Authenticate(cred.Password); err != nil {
			http.Error(w, fmt.Sprintf("invalid credentials"), http.StatusUnauthorized)
		}
		sessions.BeginSession(ctx.signingKey, ctx.sessionStore, &SessionState{}, w)
		respond(w, user, http.StatusCreated)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ctx *Context) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {

}