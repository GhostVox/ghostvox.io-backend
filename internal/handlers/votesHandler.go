package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/GhostVox/ghostvox.io-backend/internal/config"
	"github.com/google/uuid"
)

type Vote struct {
	PollId   string       `json:"poll_id"`
	OptionId string       `json:"option_id"`
	UserId   string       `json:"user_id"`
	Poll     PollResponse `json:"poll"`
}

type voteHandler struct {
	cfg *config.APIConfig
}

func NewVoteHandler(db *config.APIConfig) *voteHandler {
	return &voteHandler{
		cfg: db,
	}
}

func (vh *voteHandler) VoteOnPoll(w http.ResponseWriter, r *http.Request) {

	// get the poll ID
	pollId := r.PathValue("pollId")
	pollUUID, err := uuid.Parse(pollId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid poll ID format", err)
		return
	}

	// parse request body into a Vote struct
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

	userUUID, err := uuid.Parse(vote.UserId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid user ID format", err)
		return
	}

	err = CreateVoteAndUpdateOptionCount(r.Context(), vh.cfg, userUUID, optionUUID, pollUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Failed to create vote", err)
		return
	}
	vote.Poll.Votes++

	respondWithJSON(w, http.StatusCreated, vote.Poll)

}
