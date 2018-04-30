package handlers

import (
	"net/http"
	"encoding/json"
	"log"
)

func respond(w http.ResponseWriter, contentType string, value interface{}, statusCode int) {
	w.Header().Add(headerContentType, contentType)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("error encoding JSON: %v", err)
	}
}