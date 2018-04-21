package handlers

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/challenges-aaronluannguyen/servers/gateway/models/users"
	"github.com/challenges-aaronluannguyen/servers/gateway/sessions"
	"path"
	"strconv"
	"time"
)

//TODO: define HTTP handler functions as described in the
//assignment description. Remember to use your handler context
//struct as the receiver on these functions so that you have
//access to things like the session store and user store.
func (ctx *Context) UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodPost:
			if err := checkJSONType(w, r); err != nil {
				return
			}
			newUser := users.NewUser{}
			if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
				http.Error(w, fmt.Sprintf("error decoding into new user: %v", err), http.StatusBadRequest)
				return
			}
			if err := newUser.Validate(); err != nil {
				http.Error(w, fmt.Sprintf("invalid data for new user: %v", err), http.StatusBadRequest)
				return
			}
			newToUser, err := newUser.ToUser()
			if err != nil {
				http.Error(w, fmt.Sprintf("error converting new user to user: %v", err), http.StatusInternalServerError)
				return
			}
			user, err := ctx.usersStore.Insert(newToUser)
			if err != nil {
				http.Error(w, fmt.Sprintf("error inserting user into database: %v", err), http.StatusInternalServerError)
				return
			}
			newSessionState := SessionState{
				time.Now(),
				user,
			}
			_, err = sessions.BeginSession(ctx.signingKey, ctx.sessionStore, newSessionState, w)
			if err != nil {
				http.Error(w, fmt.Sprintf("error beginning session: %v", err), http.StatusInternalServerError)
				return
			}
			respond(w, contentTypeJSON, user, http.StatusCreated)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
	}
}

func (ctx *Context) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			userID, err := strconv.Atoi(path.Base(r.URL.Path))
			if err != nil {
				http.Error(w, fmt.Sprintf("error getting user id: %v", err), http.StatusBadRequest)
				return
			}
			user, err := ctx.usersStore.GetByID(int64(userID))
			if err != nil {
				http.Error(w, "no user found", http.StatusNotFound)
				return
			}
			respond(w, contentTypeJSON, user, http.StatusOK)

		case http.MethodPatch:
			userIDAsString := path.Base(r.URL.Path)
			userID, err := strconv.Atoi(userIDAsString)
			if err != nil {
				http.Error(w, fmt.Sprintf("error getting user id: %v", err), http.StatusBadRequest)
				return
			}
			currentState := &SessionState{}
			_, err = sessions.GetState(r, ctx.signingKey, ctx.sessionStore, currentState)
			if err != nil {
				http.Error(w, fmt.Sprintf("error getting session state: %v", err), http.StatusInternalServerError)
				return
			}
			if userIDAsString != "me" || int64(userID) != currentState.user.ID {
				http.Error(w,"not authorized, action is forbidden", http.StatusForbidden)
				return
			}
			if err = checkJSONType(w, r); err != nil {
				http.Error(w, fmt.Sprintf("error: request body must contain json: %v", err), http.StatusUnsupportedMediaType)
				return
			}
			userUpdate := users.Updates{}
			respond(w, contentTypeJSON, userUpdate, http.StatusOK)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
	}
}

func (ctx *Context) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodPost:
			if err := checkJSONType(w, r); err != nil {
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
			_, err = sessions.BeginSession(ctx.signingKey, ctx.sessionStore, &SessionState{}, w)
			if err != nil {
				http.Error(w, fmt.Sprintf("error beginning session: %v", err), http.StatusInternalServerError)
				return
			}
			respond(w, contentTypeJSON, user, http.StatusCreated)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
	}
}

func (ctx *Context) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodDelete:
			if path.Base(r.URL.Path) != "mine" {
				http.Error(w, "not authorized to delete current session", http.StatusForbidden)
				return
			}
			_, err := sessions.EndSession(r, ctx.signingKey, ctx.sessionStore)
			if err != nil {
				http.Error(w, fmt.Sprintf("error ending session: %v", err), http.StatusInternalServerError)
				return
			}
			w.Write([]byte("signed out"))

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
	}
}