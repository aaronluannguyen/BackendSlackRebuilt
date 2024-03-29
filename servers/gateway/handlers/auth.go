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
	"github.com/challenges-aaronluannguyen/servers/gateway/indexes"
	"strings"
	"github.com/nbutton23/zxcvbn-go"
)

const headerCORS = "Access-Control-Allow-Origin"
const headerContentType = "Content-Type"
const contentTypeJSON = "application/json"

//TODO: define HTTP handler functions as described in the
//assignment description. Remember to use your handler context
//struct as the receiver on these functions so that you have
//access to things like the session store and user store.
func (ctx *Context) UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			currentState := &SessionState{}
			_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, currentState)
			if err != nil {
				http.Error(w, fmt.Sprintf("unauthenticated user error: %v", err), http.StatusUnauthorized)
				return
			}
			query := r.URL.Query().Get("q")
			if len(query) < 1 {
				http.Error(w, "query string parameter required", http.StatusBadRequest)
				return
			}

			topTwentyUsers := ctx.Trie.Find(query, 20)
			var sortedTopUsers *[]*users.User
			if len(topTwentyUsers) > 0 {
				sortedTopUsers, err = ctx.UsersStore.SortTopTwentyUsersByUsername(topTwentyUsers)
				if err != nil {
					http.Error(w, fmt.Sprintf("error retrieving top users: %v", err), http.StatusBadRequest)
					return
				}
			}
			respond(w, contentTypeJSON, sortedTopUsers, http.StatusCreated)

		case http.MethodPost:
			if err := checkJSONType(w, r); err != nil {
				http.Error(w, fmt.Sprintf("error: request body must contain json: %v", err), http.StatusUnsupportedMediaType)
				return
			}
			newUser := &users.NewUser{}
			if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
				http.Error(w, fmt.Sprintf("error decoding into new user: %v", err), http.StatusBadRequest)
				return
			}
			if err := newUser.Validate(); err != nil {
				http.Error(w, fmt.Sprintf("invalid data for new user: %v", err), http.StatusBadRequest)
				return
			}
			passwordScore := zxcvbn.PasswordStrength(newUser.Password, nil)
			if passwordScore.Score <= 2 {
				http.Error(w, fmt.Sprintf("password strength is not safely unguessable"), http.StatusBadRequest)
				return
			}
			newToUser, err := newUser.ToUser()
			if err != nil {
				http.Error(w, fmt.Sprintf("error converting new user to user: %v", err), http.StatusInternalServerError)
				return
			}
			user, err := ctx.UsersStore.Insert(newToUser)
			if err != nil {
				http.Error(w, fmt.Sprintf("error inserting user into database: %v", err), http.StatusInternalServerError)
				return
			}
			users.AddUserToTrie(ctx.Trie, newToUser)
			newSessionState := &SessionState{
				time.Now(),
				user,
			}
			_, err = sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, newSessionState, w)
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
			user, err := ctx.UsersStore.GetByID(int64(userID))
			if err != nil {
				http.Error(w, "no user found", http.StatusNotFound)
				return
			}
			respond(w, contentTypeJSON, user, http.StatusOK)

		case http.MethodPatch:
			currentState := &SessionState{}
			_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, currentState)
			if err != nil {
				http.Error(w, fmt.Sprintf("error getting session state: %v", err), http.StatusUnauthorized)
				return
			}
			var userID int64
			userIDAsString := path.Base(r.URL.Path)
			if userIDAsString != "me" {
				userID, err = strconv.ParseInt(userIDAsString, 10, 64)
				if err != nil {
					http.Error(w, fmt.Sprintf("error getting user id: %v", err), http.StatusBadRequest)
					return
				}
				if int64(userID) != currentState.User.ID {
					http.Error(w,"not authorized, action is forbidden", http.StatusForbidden)
					return
				}
			} else {
				userID = currentState.User.ID
			}
			if err := checkJSONType(w, r); err != nil {
				http.Error(w, fmt.Sprintf("error: request body must contain json: %v", err), http.StatusUnsupportedMediaType)
				return
			}
			userUpdate := &users.Updates{}
			if err := json.NewDecoder(r.Body).Decode(userUpdate); err != nil {
				http.Error(w, fmt.Sprintf("error decoding into credentials: %v", err), http.StatusBadRequest)
				return
			}
			oldUser := currentState.User
			if err := currentState.User.ApplyUpdates(userUpdate); err != nil {
				http.Error(w, fmt.Sprintf("error with updating user: %s", err), http.StatusBadRequest)
				return
			}
			user, err := ctx.UsersStore.Update(userID, userUpdate)
			if err != nil {
				http.Error(w, fmt.Sprintf("error updating user: %s", err), http.StatusInternalServerError)
			}
			TrieHandleUserUpdate(ctx.Trie, oldUser, user)
			respond(w, contentTypeJSON, user, http.StatusOK)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
	}
}

func (ctx *Context) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodPost:
			if err := checkJSONType(w, r); err != nil {
				http.Error(w, fmt.Sprintf("request body must be json type error: %v", err), http.StatusBadRequest)
				return
			}
			cred := &users.Credentials{}
			if err := json.NewDecoder(r.Body).Decode(cred); err != nil {
				http.Error(w, fmt.Sprintf("error decoding into credentials: %v", err), http.StatusBadRequest)
				return
			}
			user, err := ctx.UsersStore.GetByEmail(cred.Email)
			if err != nil {
				http.Error(w, fmt.Sprintf("invalid credentials"), http.StatusUnauthorized)
				return
			}
			if err := user.Authenticate(cred.Password); err != nil {
				http.Error(w, fmt.Sprintf("invalid credentials"), http.StatusUnauthorized)
				return
			}
			newSessionState := &SessionState{
				time.Now(),
				user,
			}
			_, err = sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, newSessionState, w)
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
			_, err := sessions.EndSession(r, ctx.SigningKey, ctx.SessionStore)
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

//TrieHandlerUserUpdate removes the user's old first and last names from the trie
//and updates the trie with the new first and last names
func TrieHandleUserUpdate(trie *indexes.Trie, oldUser *users.User, updatedUser *users.User) {
	oldFirstName := strings.Split(strings.ToLower(oldUser.FirstName), " ")
	RemoveOldFirstOrLast(trie, oldFirstName, oldUser.ID)

	oldLastName := strings.Split(strings.ToLower(oldUser.LastName), " ")
	RemoveOldFirstOrLast(trie, oldLastName, oldUser.ID)

	newFirstName := strings.Split(strings.ToLower(updatedUser.FirstName), " ")
	AddUpdatedFirstOrLast(trie, newFirstName, updatedUser.ID)

	newLastName := strings.Split(strings.ToLower(updatedUser.LastName), " ")
	AddUpdatedFirstOrLast(trie, newLastName, updatedUser.ID)
}

//RemoveOldFirstOrLast removes either the old first or last names from the trie
func RemoveOldFirstOrLast(trie *indexes.Trie, names []string, id int64) {
	for _, name := range names {
		trimName := strings.TrimSpace(name)
		trie.Remove(trimName, id)
	}
}

//AddUpdatedFirstOrLast adds the new updated first or last names to the trie
func AddUpdatedFirstOrLast(trie *indexes.Trie, names []string, id int64) {
	for _, name := range names {
		trimName := strings.TrimSpace(name)
		trie.Add(trimName, id)
	}
}