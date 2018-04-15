package sessions

import (
	"errors"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "
const schemeBearerNoSpace = "Bearer"

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	//TODO:
	//- create a new SessionID
	//- save the sessionState to the store
	//- add a header to the ResponseWriter that looks like this:
	//    "Authorization: Bearer <sessionID>"
	//  where "<sessionID>" is replaced with the newly-created SessionID
	//  (note the constants declared for you above, which will help you avoid typos)

	newSessionID, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, nil
	}
	err = store.Save(newSessionID, sessionState)
	if err != nil {
		return InvalidSessionID, err
	}
	w.Header().Add(headerAuthorization, schemeBearer + string(newSessionID))
	return newSessionID, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	//TODO: get the value of the Authorization header,
	//or the "auth" query string parameter if no Authorization header is present,
	//and validate it. If it's valid, return the SessionID. If not
	//return the validation error.

	authVals := strings.Split(r.Header.Get(headerAuthorization), " ")
	isValid := checkAuthIsValid(authVals)
	switch {
		case isValid: return ValidateID(authVals[1], signingKey)
		case !isValid:
			authParamVals := strings.Split(r.URL.Query().Get(paramAuthorization), " ")
			isParamValid := checkAuthIsValid(authParamVals)
			if isParamValid {
				return ValidateID(authParamVals[1], signingKey)
			}
			return InvalidSessionID, ErrInvalidScheme
	}
	return InvalidSessionID, ErrNoSessionID
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	//TODO: get the SessionID from the request, and get the data
	//associated with that SessionID from the store.
	return InvalidSessionID, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	//TODO: get the SessionID from the request, and delete the
	//data associated with it in the store.
	return InvalidSessionID, nil
}


func checkAuthIsValid(values []string) bool {
	switch {
	case len(values) <= 1:
		return false
	case len(values) == 2:
		if values[0] == schemeBearerNoSpace {
			return true
		}
	}
	return false
}