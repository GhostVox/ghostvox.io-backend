package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/GhostVox/ghostvox.io-backend/internal/database"
	"github.com/google/uuid"
)

type Vote struct {
	PollId   string `json:"poll_id"`
	OptionId string `json:"option_id"`
	UserId   string `json:"user_id"`
}

type voteHandler struct {
	db *database.Queries
}

func NewVoteHandler(db *database.Queries) *voteHandler {
	return &voteHandler{
		db: db,
	}
}
func (vh *voteHandler) GetVotesByPoll(w http.ResponseWriter, r *http.Request) {

	pollId := r.PathValue("pollId")
	pollUUID, err := uuid.Parse(pollId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Missing poll id path value", err)
		return
	}
	votes, err := vh.db.GetVotesByPollID(context.Background(), pollUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), "No votes found for poll", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Failed to retrieve votes", err)
		return
	}

	respondWithJSON(w, http.StatusOK, votes)

}

func (vh *voteHandler) CreateVote(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("pollId")
	pollUUID, err := uuid.Parse(pollId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid poll ID format", err)
		return
	}
	var vote Vote
	err = json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid request body", err)
		return
	}

	optionUUID, err := uuid.Parse(vote.OptionId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid option ID format", err)
		return
	}

	voteRecord, err := vh.db.CreateVote(r.Context(), database.CreateVoteParams{
		PollID:   pollUUID,
		OptionID: optionUUID,
		UserID:   vote.UserId,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), "Poll or option not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Failed to create vote", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, voteRecord)

}

func (vh *voteHandler) DeleteVote(w http.ResponseWriter, r *http.Request) {

	voteId := r.PathValue("voteId")
	voteUUID, err := uuid.Parse(voteId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid vote ID format", err)
		return
	}

	err = vh.db.DeleteVoteByID(r.Context(), voteUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), "Vote not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Failed to delete vote", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
