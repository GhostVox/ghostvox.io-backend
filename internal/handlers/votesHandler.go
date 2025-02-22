package handlers

import (
	"net/http"

	"github.com/GhostVox/ghostvox.io-backend/internal/database"
)

type voteHandler struct {
	db *database.Queries
}

func NewVoteHandler(db *database.Queries) *voteHandler {
	return &voteHandler{
		db: db,
	}
}

func (vh *voteHandler) GetVote(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to create an option
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
	return
}

func (vh *voteHandler) GetVotes(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to create an option
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
	return
}

func (vh *voteHandler) CreateVote(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to create an option
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
	return
}

func (vh *voteHandler) UpdateVote(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to create an option
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
	return
}

func (vh *voteHandler) DeleteVote(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to create an option
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
	return
}
