package handlers

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/challenges-aaronluannguyen/servers/gateway/models/users"
	"github.com/challenges-aaronluannguyen/servers/gateway/sessions"
	"time"
	"bytes"
	"encoding/json"
	"reflect"
)

func createNormUser() *users.User {
	nu := &users.NewUser{
		Email: "test1@uw.edu",
		Password: "123123",
		PasswordConf: "123123",
		UserName: "tester1",
		FirstName: "first",
		LastName: "last",
	}
	user, _ := nu.ToUser()
	return user
}

func TestUsersHandler(t *testing.T) {
	cases := []struct {
		method string
		codeExpected int
	} {
		{
			"GET",
			http.StatusMethodNotAllowed,
		},
		{
			"PUT",
			http.StatusMethodNotAllowed,
		},
		{
			"PATCH",
			http.StatusMethodNotAllowed,
		},
		{
			"DELETE",
			http.StatusMethodNotAllowed,
		},
	}
	for _, c := range cases {
		req, err := http.NewRequest(c.method, "/users", nil)
		if err != nil {
			t.Errorf("error sending requests.")
		}
		req.Header.Set("Content-Type", contentTypeJSON)
		rr := httptest.NewRecorder()
		ctx := Context{}
		handler := http.HandlerFunc(ctx.UsersHandler)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != c.codeExpected {
			t.Errorf("handler returned wrong status code: got %v but expect %v", status, http.StatusMethodNotAllowed)
		}
	}

	tests := []struct {
		name string
		err bool
		method string
		newU string
		sessStore *sessions.MemStore
		userStore *users.MockStore
		contentType string
	} {
		{
			"Check Valid Json",
			true,
			"POST",
			"",
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(true, nil),
			"",
		},
		{
			"Check valid new user",
			false,
			"POST",
			`{"email": "test1@uw.edu", "password": "123123", "passwordConf": "123123", "userName": "tester1", "firstName": "first", "lastName": "last"}`,
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(false, &users.User{ID: int64(1), UserName: "tester1", FirstName: "first", LastName: "last", PhotoURL: "https://www.gravatar.com/avatar/ae0f326d14556306ff28ea5e485796e9"}),
			"application/json",
		},
		{
			"Check invalid new user",
			false,
			"POST",
			`{"email": "test1@uw.edu", "password": "123123", "passwordConf": "123123", "userName": "test er1", "firstName": "first", "lastName": "last"}`,
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(true, &users.User{ID: int64(1), UserName: "tester1", FirstName: "first", LastName: "last", PhotoURL: "https://www.gravatar.com/avatar/ae0f326d14556306ff28ea5e485796e9"}),
			"application/json",
		},
		{
			"Check decode err",
			false,
			"POST",
			`}{`,
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(true, &users.User{ID: int64(1), UserName: "tester1", FirstName: "first", LastName: "last", PhotoURL: "https://www.gravatar.com/avatar/ae0f326d14556306ff28ea5e485796e9"}),
			"application/json",
		},
		{
			"Check session insert error",
			false,
			"POST",
			`{"email": "test1@uw.edu", "password": "123123", "passwordConf": "123123", "userName": "tester1", "firstName": "first", "lastName": "last"}`,
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(true, &users.User{ID: int64(1), UserName: "tester1", FirstName: "first", LastName: "last", PhotoURL: "https://www.gravatar.com/avatar/ae0f326d14556306ff28ea5e485796e9"}),
			"application/json",
		},
	}

	for _, tt := range tests {
		byteJSON := []byte(tt.newU)
		queryJSON := bytes.NewBuffer(byteJSON)
		req, err := http.NewRequest(tt.method, "/v1/users", queryJSON)
		if err != nil {
			t.Errorf("error sending requests. case: %s", tt.name)
		}

		req.Header.Set("Content-Type", tt.contentType)
		rr := httptest.NewRecorder()
		ctx := Context{"key", tt.sessStore, tt.userStore}
		ctx.UsersHandler(rr, req)
		if rr.Result().StatusCode < 300 && tt.err {
			t.Errorf("case %s: expected error but didn't get one", tt.name)
		}
		if rr.Result().StatusCode == http.StatusOK || rr.Result().StatusCode == http.StatusCreated {
			testU := &users.User{}
			if err := json.Unmarshal(rr.Body.Bytes(), testU); err != nil {
				t.Errorf("case %s: error unmarshalling json", tt.name)
			}
			if !reflect.DeepEqual(tt.userStore.User, testU) {
				t.Errorf("case %s: result not equal to expected result", tt.name)
			}
		}
	}
}

func getSessionID (key string) sessions.SessionID {
	id, _ := sessions.NewSessionID(key)
	return id
}

func TestSpecificUserHandler(t *testing.T) {
	cases := []struct {
		method string
		codeExpected int
	} {
		{
			"PUT",
			http.StatusMethodNotAllowed,
		},
		{
			"POST",
			http.StatusMethodNotAllowed,
		},
		{
			"DELETE",
			http.StatusMethodNotAllowed,
		},
	}
	for _, c := range cases {
		req, err := http.NewRequest(c.method, "/users", nil)
		if err != nil {
			t.Errorf("error sending requests.")
		}
		req.Header.Set("Content-Type", contentTypeJSON)
		rr := httptest.NewRecorder()
		ctx := Context{}
		handler := http.HandlerFunc(ctx.SpecificUserHandler)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != c.codeExpected {
			t.Errorf("handler returned wrong status code: got %v but expect %v", status, http.StatusMethodNotAllowed)
		}
	}

	tests := []struct {
		name string
		err bool
		method string
		newU string
		sessStore *sessions.MemStore
		userStore *users.MockStore
		contentType string
		userID string
	} {
		{
			"Check valid Get",
			false,
			"GET",
			`{"email": "test1@uw.edu", "password": "123123", "passwordConf": "123123", "userName": "tester1", "firstName": "first", "lastName": "last"}`,
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(false, &users.User{ID: int64(1), UserName: "tester1", FirstName: "first", LastName: "last", PhotoURL: "https://www.gravatar.com/avatar/ae0f326d14556306ff28ea5e485796e9"}),
			"",
			"1",
		},
		{
			"Check invalid user id",
			true,
			"GET",
			`{"email": "test1@uw.edu", "password": "123123"}`,
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(true, &users.User{ID: int64(1), UserName: "tester1", FirstName: "first", LastName: "last", PhotoURL: "https://www.gravatar.com/avatar/ae0f326d14556306ff28ea5e485796e9"}),
			"application/json",
			"12qwe",
		},
		{
			"Check User Store Get Err",
			false,
			"GET",
			`{"email": "test1@uw.edu", "password": "123123", "passwordConf": "123123", "userName": "tester1", "firstName": "first", "lastName": "last"}`,
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(true, &users.User{ID: int64(1), UserName: "tester1", FirstName: "first", LastName: "last", PhotoURL: "https://www.gravatar.com/avatar/ae0f326d14556306ff28ea5e485796e9"}),
			"application/json",
			"1",
		},
		{
			"Session state get err",
			true,
			"PATCH",
			`{"email": "test1@uw.edu", "password": "123123", "passwordConf": "123123", "userName": "tester1", "firstName": "first", "lastName": "last"}`,
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(true, &users.User{ID: int64(1), UserName: "tester1", FirstName: "first", LastName: "last", PhotoURL: "https://www.gravatar.com/avatar/ae0f326d14556306ff28ea5e485796e9"}),
			"application/json",
			"1",
		},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, "/v1/users/" + tt.userID, nil)
		if err != nil {
			t.Errorf("error sending requests. case: %s", tt.name)
		}

		req.Header.Set("Content-Type", tt.contentType)
		rr := httptest.NewRecorder()
		ctx := Context{"key", tt.sessStore, tt.userStore}
		ctx.SpecificUserHandler(rr, req)
		if rr.Result().StatusCode < 300 && tt.err {
			t.Errorf("case %s: expected error but didn't get one", tt.name)
		}
		if rr.Result().StatusCode == http.StatusOK || rr.Result().StatusCode == http.StatusCreated {
			testU := &users.User{}
			if err := json.Unmarshal(rr.Body.Bytes(), testU); err != nil {
				t.Errorf("case %s: error unmarshalling json", tt.name)
			}
			if !reflect.DeepEqual(tt.userStore.User, testU) {
				t.Errorf("case %s: result not equal to expected result", tt.name)
			}
		}
	}

	testsPatch := []struct {
		name string
		err bool
		method string
		newU string
		sessStore *sessions.MemStore
		userStore *users.MockStore
		contentType string
		userID string
		sessionID sessions.SessionID
	} {
		{
			"Check valid Patch",
			false,
			"PATCH",
			`{"firstName": "first", "lastName": "last"}`,
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(false, &users.User{ID: int64(1), UserName: "tester1", FirstName: "first", LastName: "last", PhotoURL: "https://www.gravatar.com/avatar/ae0f326d14556306ff28ea5e485796e9"}),
			"application/json",
			"1",
			getSessionID("key"),
		},
		{
			"Check valid Patch with me",
			false,
			"PATCH",
			`{"firstName": "first", "lastName": "last"}`,
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(false, &users.User{ID: int64(1), UserName: "tester1", FirstName: "first", LastName: "last", PhotoURL: "https://www.gravatar.com/avatar/ae0f326d14556306ff28ea5e485796e9"}),
			"application/json",
			"me",
			getSessionID("key"),
		},
		{
			"Check non existent json content type",
			true,
			"PATCH",
			`{"firstName": "first", "lastName": "last"}`,
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(false, &users.User{ID: int64(1), UserName: "tester1", FirstName: "first", LastName: "last", PhotoURL: "https://www.gravatar.com/avatar/ae0f326d14556306ff28ea5e485796e9"}),
			"",
			"me",
			getSessionID("key"),
		},
		{
			"Check invalid id containing letters",
			false,
			"PATCH",
			`{"firstName": "first", "lastName": "last"}`,
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(false, &users.User{ID: int64(1), UserName: "tester1", FirstName: "first", LastName: "last", PhotoURL: "https://www.gravatar.com/avatar/ae0f326d14556306ff28ea5e485796e9"}),
			"application/json",
			"me123",
			getSessionID("key"),
		},
		{
			"Check invalid patch json",
			false,
			"PATCH",
			"{,,",
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(false, &users.User{ID: int64(1), UserName: "tester1", FirstName: "first", LastName: "last", PhotoURL: "https://www.gravatar.com/avatar/ae0f326d14556306ff28ea5e485796e9"}),
			"application/json",
			"me",
			getSessionID("key"),
		},
		{
			"Check user store update err",
			true,
			"PATCH",
			`{"firstName": "first", "lastName": "last"}`,
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(true, &users.User{ID: int64(1), UserName: "tester1", FirstName: "first", LastName: "last", PhotoURL: "https://www.gravatar.com/avatar/ae0f326d14556306ff28ea5e485796e9"}),
			"application/json",
			"me",
			getSessionID("key"),
		},
		{
			"Check apply update err",
			true,
			"PATCH",
			`{"firstName": "", "lastName": ""}`,
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(true, &users.User{ID: int64(1), UserName: "tester1", FirstName: "first", LastName: "last", PhotoURL: "https://www.gravatar.com/avatar/ae0f326d14556306ff28ea5e485796e9"}),
			"application/json",
			"me",
			getSessionID("key"),
		},
	}

	for _, tt := range testsPatch {
		byteJSON := []byte(tt.newU)
		queryJSON := bytes.NewBuffer(byteJSON)
		req, err := http.NewRequest(tt.method, "/v1/users/" + tt.userID, queryJSON)
		if err != nil {
			t.Errorf("error sending requests. case: %s", tt.name)
		}
		req.Header.Set("Authorization", "Bearer " + tt.sessionID.String())
		req.Header.Set("Content-Type", tt.contentType)
		rr := httptest.NewRecorder()
		currState := &SessionState{time.Now(), tt.userStore.User}
		tt.sessStore.Save(tt.sessionID, currState)
		ctx := Context{"key", tt.sessStore, tt.userStore}
		ctx.SpecificUserHandler(rr, req)
		if rr.Result().StatusCode < 300 && tt.err {
			t.Errorf("case %s: expected error but didn't get one", tt.name)
		}
		if rr.Result().StatusCode == http.StatusOK || rr.Result().StatusCode == http.StatusCreated {
			testU := &users.User{}
			if err := json.Unmarshal(rr.Body.Bytes(), testU); err != nil {
				t.Errorf("case %s: error unmarshalling json", tt.name)
			}
			if !reflect.DeepEqual(tt.userStore.User, testU) {
				t.Errorf("case %s: result not equal to expected result", tt.name)
			}
		}
	}
}

func TestSessionsHandler(t *testing.T) {
	cases := []struct {
		method string
		codeExpected int
	} {
		{
			"GET",
			http.StatusMethodNotAllowed,
		},
		{
			"PUT",
			http.StatusMethodNotAllowed,
		},
		{
			"PATCH",
			http.StatusMethodNotAllowed,
		},
		{
			"DELETE",
			http.StatusMethodNotAllowed,
		},
	}
	for _, c := range cases {
		req, err := http.NewRequest(c.method, "/users", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", contentTypeJSON)
		rr := httptest.NewRecorder()
		ctx := Context{}
		handler := http.HandlerFunc(ctx.SessionsHandler)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != c.codeExpected {
			t.Errorf("handler returned wrong status code: got %v but expect %v", status, http.StatusMethodNotAllowed)
		}
	}

	tests := []struct {
		name string
		err bool
		method string
		newU string
		sessStore *sessions.MemStore
		userStore *users.MockStore
		contentType string
	} {
		{
			"Check Valid Json",
			true,
			"POST",
			"",
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(true, nil),
			"",
		},
		{
			"Check Valid Credentials Struct",
			false,
			"POST",
			`{"email": "test1@uw.edu", "password": "123123"}`,
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(false, &users.User{ID: int64(1), UserName: "tester1", FirstName: "first", LastName: "last", PhotoURL: "https://www.gravatar.com/avatar/ae0f326d14556306ff28ea5e485796e9"}),
			"application/json",
		},
		{
			"Check invalid Credentials Struct",
			true,
			"POST",
			"",
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(true, &users.User{ID: int64(1), UserName: "tester1", FirstName: "first", LastName: "last", PhotoURL: "https://www.gravatar.com/avatar/ae0f326d14556306ff28ea5e485796e9"}),
			"application/json",
		},
		{
			"Check invalid store get by email",
			true,
			"POST",
			`{"email": "test1@uw.edu", "password": "123123"}`,
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(true, nil),
			"application/json",
		},
		{
			"Check valid login",
			false,
			"POST",
			`{"email": "test1@uw.edu", "password": "123123"}`,
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(false, createNormUser()),
			"application/json",
		},
	}

	for _, tt := range tests {
		byteJSON := []byte(tt.newU)
		queryJSON := bytes.NewBuffer(byteJSON)
		req, err := http.NewRequest(tt.method, "/v1/sessions", queryJSON)
		if err != nil {
			t.Errorf("error sending requests. case: %s", tt.name)
		}

		req.Header.Set("Content-Type", tt.contentType)
		rr := httptest.NewRecorder()
		ctx := Context{"key", tt.sessStore, tt.userStore}
		ctx.SessionsHandler(rr, req)
		if rr.Result().StatusCode < 300 && tt.err {
			t.Errorf("case %s: expected error but didn't get one", tt.name)
		}
		if rr.Result().StatusCode == http.StatusOK || rr.Result().StatusCode == http.StatusCreated {
			testU := &users.Credentials{}
			if err := json.Unmarshal(rr.Body.Bytes(), testU); err != nil {
				t.Errorf("case %s: error unmarshalling json", tt.name)
			}
		}
	}
}

func TestSpecficSessionHandler(t *testing.T) {
	cases := []struct {
		method string
		codeExpected int
	} {
		{
			"GET",
			http.StatusMethodNotAllowed,
		},
		{
			"PUT",
			http.StatusMethodNotAllowed,
		},
		{
			"PATCH",
			http.StatusMethodNotAllowed,
		},
		{
			"POST",
			http.StatusMethodNotAllowed,
		},
	}
	for _, c := range cases {
		req, err := http.NewRequest(c.method, "/users", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", contentTypeJSON)
		rr := httptest.NewRecorder()
		ctx := Context{}
		handler := http.HandlerFunc(ctx.SpecificSessionHandler)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != c.codeExpected {
			t.Errorf("handler returned wrong status code: got %v but expect %v", status, http.StatusMethodNotAllowed)
		}
	}

	tests := []struct {
		name string
		err bool
		method string
		newU string
		sessStore *sessions.MemStore
		userStore *users.MockStore
		contentType string
		query string
	} {
		{
			"Check non mine request",
			true,
			"DELETE",
			"",
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(true, nil),
			"",
			"",
		},
		{
			"Check end session err",
			true,
			"DELETE",
			"",
			sessions.NewMemStore(time.Hour, time.Minute),
			users.NewMockStore(true, nil),
			"",
			"mine",
		},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, "/v1/sessions/" + tt.query, nil)
		if err != nil {
			t.Errorf("error sending requests. case: %s", tt.name)
		}

		req.Header.Set("Content-Type", tt.contentType)
		rr := httptest.NewRecorder()
		ctx := Context{"key", tt.sessStore, tt.userStore}
		ctx.SpecificSessionHandler(rr, req)
		if rr.Result().StatusCode < 300 && tt.err {
			t.Errorf("case %s: expected error but didn't get one", tt.name)
		}
		if rr.Result().StatusCode == http.StatusOK || rr.Result().StatusCode == http.StatusCreated {
			testU := &users.Credentials{}
			if err := json.Unmarshal(rr.Body.Bytes(), testU); err != nil {
				t.Errorf("case %s: error unmarshalling json", tt.name)
			}
			if !reflect.DeepEqual(tt.userStore.User, testU) {
				t.Errorf("case %s: result not equal to expected result", tt.name)
			}
		}
	}
}