package handlers

import (
	c "context"
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
	polls, err := h.db.GetAllPolls(c.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusNotFound, err)
			return
		} else {
			chooseError(w, http.StatusInternalServerError, err)
			return
		}
	}

	respondWithJSON(w, http.StatusOK, polls)
	return
}

func (h *pollHandler) GetPoll(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("pollId")
	if id == "" {
		chooseError(w, http.StatusBadRequest, errors.New("pollId is Required"))
		return
	}

	pollUUID, err := uuid.Parse(id)
	if err != nil {
		chooseError(w, http.StatusBadRequest, errors.New("Invalid pollId"))
		return
	}

	poll, err := h.db.GetPoll(c.Background(), pollUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusNotFound, err)
			return
		} else {
			chooseError(w, http.StatusInternalServerError, err)
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
		chooseError(w, http.StatusInternalServerError, err)
		return
	}

	userUUID, err := uuid.Parse(newPoll.UserID)
	if err != nil {
		chooseError(w, http.StatusBadRequest, errors.New("Invalid UserID"))
		return
	}

	expiresAt, err := time.Parse(time.RFC3339, newPoll.ExpiresAt)
	if err != nil {
		chooseError(w, http.StatusBadRequest, errors.New("Invalid ExpiresAt"))
		return
	}
	if newPoll.ExpiresAt == "" {
		expiresAt = time.Now().UTC().Add(time.Duration(24 * time.Hour))
	}

	pollRecord, err := h.db.CreatePoll(c.Background(), database.CreatePollParams{
		UserID:      userUUID,
		Description: newPoll.Description,
		Title:       newPoll.Description,
		ExpiresAt:   expiresAt,
		Status:      database.PollStatus(newPoll.Status),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusNotFound, err)
			return
		} else {
			chooseError(w, http.StatusInternalServerError, err)
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
		chooseError(w, http.StatusInternalServerError, err)
		return
	}

	userUUID, err := uuid.Parse(newPoll.UserID)
	if err != nil {
		chooseError(w, http.StatusBadRequest, errors.New("Invalid UserID"))
		return
	}

	pollUUID, err := uuid.Parse(pollId)
	if err != nil {
		chooseError(w, http.StatusBadRequest, errors.New("Invalid PollID"))
		return
	}
	expiresAt, err := time.Parse(time.RFC3339, newPoll.ExpiresAt)
	if err != nil {
		chooseError(w, http.StatusBadRequest, errors.New("Invalid ExpiresAt expect format to be 2023-10-05T15:04:05Z07:00"))
		return
	}

	pollRecord, err := h.db.UpdatePoll(c.Background(), database.UpdatePollParams{
		ID:          pollUUID,
		UserID:      userUUID,
		Description: newPoll.Description,
		Title:       newPoll.Description,
		ExpiresAt:   expiresAt,
		Status:      database.PollStatus(newPoll.Status),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusNotFound, err)
			return
		} else {
			chooseError(w, http.StatusInternalServerError, err)
			return
		}
	}

	respondWithJSON(w, http.StatusOK, pollRecord)
	return
}

func (h *pollHandler) DeletePoll(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("pollId")
	if pollId == "" {
		chooseError(w, http.StatusBadRequest, errors.New("PollID is required"))
		return
	}

	pollUUID, err := uuid.Parse(pollId)
	if err != nil {
		chooseError(w, http.StatusBadRequest, errors.New("Invalid PollID"))
		return
	}

	err = h.db.DeletePoll(c.Background(), pollUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusNotFound, err)
			return
		} else {
			chooseError(w, http.StatusInternalServerError, err)
			return
		}
	}
	respondWithJSON(w, http.StatusNoContent, nil)
	return
}
