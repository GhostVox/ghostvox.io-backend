package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/GhostVox/ghostvox.io-backend/internal/database"
	"github.com/google/uuid"
)

type poll struct {
	UserID      string `json:"userId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ExpiresAt   string `json:"expiresAt"`
	Status      string `json:"status"`
}

type pollHandler struct {
	db *database.Queries
}

func NewPollHandler(db *database.Queries) *pollHandler {
	return &pollHandler{
		db: db,
	}
}

// Polls route
func (h *pollHandler) GetAllPolls(w http.ResponseWriter, r *http.Request) {
	polls, err := h.db.GetAllPolls(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), "No polls found", err)
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
			return
		}
	}

	respondWithJSON(w, http.StatusOK, polls)
	return
}

func (h *pollHandler) GetPoll(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("pollId")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "poll id missing from pathvalue", errors.New("pollId is required"))
		return
	}

	pollUUID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid poll id", err)
		return
	}

	poll, err := h.db.GetPoll(r.Context(), pollUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), "No poll found", err)
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
			return
		}
	}

	respondWithJSON(w, http.StatusOK, poll)
	return
}

func (h *pollHandler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	newPoll := poll{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&newPoll); err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
		return
	}

	expiresAt, err := time.Parse(time.RFC3339, newPoll.ExpiresAt)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid ExpiresAt", err)
		return
	}
	if newPoll.ExpiresAt == "" {
		expiresAt = time.Now().UTC().Add(time.Duration(24 * time.Hour))
	}

	pollRecord, err := h.db.CreatePoll(r.Context(), database.CreatePollParams{
		UserID:      newPoll.UserID,
		Description: newPoll.Description,
		Title:       newPoll.Title,
		ExpiresAt:   expiresAt,
		Status:      database.PollStatus(newPoll.Status),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), "No poll found", err)
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
			return
		}
	}

	respondWithJSON(w, http.StatusCreated, pollRecord)
	return
}

func (h *pollHandler) UpdatePoll(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("pollId")
	newPoll := poll{}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&newPoll)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
		return
	}

	pollUUID, err := uuid.Parse(pollId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid poll id in pathvalue", err)
		return
	}
	expiresAt, err := time.Parse(time.RFC3339, newPoll.ExpiresAt)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid ExpiresAt expect format to be 2023-10-05T15:04:05Z07:00", err)
		return
	}

	pollRecord, err := h.db.UpdatePoll(r.Context(), database.UpdatePollParams{
		ID:          pollUUID,
		UserID:      newPoll.UserID,
		Description: newPoll.Description,
		Title:       newPoll.Title,
		ExpiresAt:   expiresAt,
		Status:      database.PollStatus(newPoll.Status),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), "Poll not found", err)
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
			return
		}
	}

	respondWithJSON(w, http.StatusOK, pollRecord)
	return
}

func (h *pollHandler) DeletePoll(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("pollId")
	if pollId == "" {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "PollID is required", nil)
		return
	}

	pollUUID, err := uuid.Parse(pollId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid poll id in pathvalue", err)
		return
	}

	err = h.db.DeletePoll(r.Context(), pollUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), "Poll not found", err)
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
			return
		}
	}
	respondWithJSON(w, http.StatusNoContent, nil)
	return
}
