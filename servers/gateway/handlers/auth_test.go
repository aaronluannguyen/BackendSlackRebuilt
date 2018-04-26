package handlers

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

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
			t.Fatal(err)
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

	//ctx := &Context{
	//	"key",
	//	sessions.NewMemStore(time.Hour, time.Hour),
	//
	//}
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
			t.Fatal(err)
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
}