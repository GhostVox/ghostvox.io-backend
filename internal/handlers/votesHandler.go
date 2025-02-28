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
		chooseError(w, http.StatusBadRequest, err)
		return
	}
	votes, err := vh.db.GetVotesByPollID(context.Background(), pollUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusNotFound, err)
			return
		}
		chooseError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, votes)

}

func (vh *voteHandler) CreateVote(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("pollId")
	pollUUID, err := uuid.Parse(pollId)
	if err != nil {
		chooseError(w, http.StatusBadRequest, err)
		return
	}
	var vote Vote
	err = json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		chooseError(w, http.StatusBadRequest, err)
		return
	}

	optionUUID, err := uuid.Parse(vote.OptionId)
	if err != nil {
		chooseError(w, http.StatusBadRequest, err)
		return
	}

	voteRecord, err := vh.db.CreateVote(r.Context(), database.CreateVoteParams{
		PollID:   pollUUID,
		OptionID: optionUUID,
		UserID:   vote.UserId,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusNotFound, err)
			return
		}
		chooseError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, voteRecord)

}

func (vh *voteHandler) DeleteVote(w http.ResponseWriter, r *http.Request) {

	voteId := r.PathValue("voteId")
	voteUUID, err := uuid.Parse(voteId)
	if err != nil {
		chooseError(w, http.StatusBadRequest, err)
		return
	}

	err = vh.db.DeleteVoteByID(r.Context(), voteUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusNotFound, err)
			return
		}
		chooseError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
