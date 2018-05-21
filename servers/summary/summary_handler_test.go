package main

import (
	"testing"
	"net/http"
	"fmt"
	"net/http/httptest"
)

func TestSummaryServer(t *testing.T) {
	cases := []struct{
		name string
		query string
		expectedStatusCode int
		expectedContentType string
	}{
		{
			"Valid URL Param",
			"url=https://github.com",
			http.StatusOK,
			contentTypeJSON,
		},
		{
			"Empty Query String",
			"",
			http.StatusBadRequest,
			"text/plain; charset=utf-8",
		},
		{
			"Invalid URL",
			"url=%20%20",
			http.StatusInternalServerError,
			"text/plain; charset=utf-8",
		},
	}

	for _, c := range cases {
		URL := fmt.Sprintf("/v1/summary?%s", c.query)
		req := httptest.NewRequest("GET", URL, nil)
		respRec := httptest.NewRecorder()
		SummaryHandler(respRec, req)

		resp := respRec.Result()

		if resp.StatusCode != c.expectedStatusCode {
			t.Errorf("case %s: incorrect status code: expected %d but got %d", c.name, c.expectedStatusCode, resp.StatusCode)
		}

		allowedOrigin := resp.Header.Get(headerCORS)
		if allowedOrigin != "*" {
			t.Errorf("case %s: incorrect incorrect CORS header: expected %s but got %s", c.name, "*", allowedOrigin)
		}
	}
}