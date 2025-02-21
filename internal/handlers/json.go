package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Printf("Error: %v", err)
	}
	if code > 499 {
		log.Printf("Responding with 5xx error: %s", err)
	}
	type ErrorResponse struct {
		Error string `json:"error"`
	}
	errRes := ErrorResponse{
		Error: msg,
	}
	respondWithJSON(w, code, errRes)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	JsonBytes, err := json.Marshal(&payload)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(code)
	_, err = w.Write(JsonBytes)
	if err != nil {
		fmt.Printf("Error writing JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
