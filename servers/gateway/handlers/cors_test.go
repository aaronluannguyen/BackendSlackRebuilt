package handlers

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"bytes"
)

func TestCORS(t *testing.T) {
	newHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	server := httptest.NewServer(WrappedCORSHandler(newHandler))
	var u bytes.Buffer
	u.WriteString(string(server.URL))
	u.WriteString("/v1/test")
	res, _ := http.Get(u.String())

	if res.Header.Get(headerCORS) != "*" {
		t.Errorf("Access-Control-Allow-Origin header not set")
	}

	if res.Header.Get(methodsCORS) != "GET, PUT, POST, PATCH, DELETE" {
		t.Errorf("Access-Control-Allow-Methods header has not been set")
	}

	if res.Header.Get(allowHeadersCORS) != "Content-Type, Authorization" {
		t.Errorf("Access-Control-Allow-Headers header has not been set")
	}

	if res.Header.Get(exposeHeadersCORS) != "Authorization" {
		t.Errorf("Access-Control-Expose-Headers header has not been set")
	}

	if res.Header.Get(maxAgeCORS) != "600" {
		t.Errorf("Access-Control-Max-Age header has not been set")
	}
}