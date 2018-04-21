package handlers

import (
	"net/http"
	"encoding/json"
	"log"
)

func respond(w http.ResponseWriter, value interface{}, statusCode int) {
	w.Header().Add(headerContentType, contentTypeJSON)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("error encoding JSON: %V", err)
	}
}