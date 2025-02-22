package handlers

import (
	"net/http"

	"github.com/GhostVox/ghostvox.io-backend/internal/database"
)

type optionHandler struct {
	db *database.Queries
}

func NewOptionHandler(db *database.Queries) *optionHandler {
	return &optionHandler{db: db}
}
func (oh *optionHandler) GetOption(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to create an option
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
	return
}

func (oh *optionHandler) GetOptions(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to create an option
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
	return
}

func (oh *optionHandler) CreateOption(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
	return
}
func (oh *optionHandler) UpdateOption(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to update an option
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
	return
}

func (oh *optionHandler) DeleteOption(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to delete an option
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
	return
}
