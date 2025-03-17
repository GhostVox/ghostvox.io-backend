package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/GhostVox/ghostvox.io-backend/internal/config"
	"github.com/GhostVox/ghostvox.io-backend/internal/database"
	"github.com/google/uuid"
)

type poll struct {
	UserID      uuid.UUID      `json:"userId"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Category    string         `json:"category"`
	ExpiresAt   int64          `json:"expiresAt"`
	Status      string         `json:"status"`
	Options     []CreateOption `json:"options"`
}

type pollHandler struct {
	cfg *config.APIConfig
}

func NewPollHandler(cfg *config.APIConfig) *pollHandler {
	return &pollHandler{
		cfg: cfg,
	}
}

// Polls route
func (h *pollHandler) GetAllPolls(w http.ResponseWriter, r *http.Request) {
	polls, err := h.cfg.Queries.GetAllPolls(r.Context())
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

	poll, err := h.cfg.Queries.GetPoll(r.Context(), pollUUID)
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

	err := CreatePollWithOptions(r.Context(), h.cfg.DB, h.cfg, newPoll)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, struct {
		msg string
	}{msg: http.StatusText(http.StatusOK)})
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
	expiresAt := time.Now().Add(time.Duration(newPoll.ExpiresAt))
	pollRecord, err := h.cfg.Queries.UpdatePoll(r.Context(), database.UpdatePollParams{
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

	err = h.cfg.Queries.DeletePoll(r.Context(), pollUUID)
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
