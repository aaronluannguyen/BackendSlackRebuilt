package handlers

import (
	"net/http"
	"fmt"
	"errors"
)

func checkJSONType(w http.ResponseWriter, r *http.Request) error {
	contentType := r.Header.Get("Content-type")
	if contentType != contentTypeJSON {
		http.Error(w, fmt.Sprintf("request body must be in JSON"), http.StatusUnsupportedMediaType)
		return errors.New("not json type")
	}
	return nil
}