package handlers

import (
	"net/http"
)

func chooseError(w http.ResponseWriter, code int, err error) {
	switch code {
	case http.StatusNotFound:
		respondWithError(w, code, "Resource not found", err)
		return
	case http.StatusUnauthorized:
		respondWithError(w, code, "Resource access denied", err)
		return
	case http.StatusBadRequest:
		respondWithError(w, code, "Bad request", err)
		return
	case http.StatusMethodNotAllowed:
		respondWithError(w, code, "Method not allowed", err)
		return
	case http.StatusNotImplemented:
		respondWithError(w, code, "Not implemented", err)
		return
	default:
		respondWithError(w, code, "Internal Server Error", err)
		return
	}
}
