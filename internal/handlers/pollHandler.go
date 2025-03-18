package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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

type ActivePollResponse struct {
	ID          uuid.UUID                        `json:"id"`
	Title       string                           `json:"title"`
	Creator     string                           `json:"creator"`
	Description string                           `json:"description"`
	Category    string                           `json:"category"`
	DaysLeft    int64                            `json:"daysLeft"`
	Options     []database.GetOptionsByPollIDRow `json:"options"`
	Votes       int64                            `json:"votes"`
	Comments    int64                            `json:"comments"`
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

func (h *pollHandler) GetAllActivePolls(w http.ResponseWriter, r *http.Request) {
	limitParam := r.URL.Query().Get("limit")
	if limitParam == "" {
		limitParam = "20"
	}

	offsetParam := r.URL.Query().Get("offset")
	if offsetParam == "" {
		offsetParam = "0"
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid limit parameter", err)
		return
	}
	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid offset parameter", err)
		return
	}

	polls, err := h.cfg.Queries.GetAllActivePollsList(r.Context(), database.GetAllActivePollsListParams{
		Limit:  int32(limit),
		Offset: int32(offset),
		Status: "Active",
	})
	fmt.Println(polls)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), "No active polls found", err)
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
			return
		}
	}
	var ResponsePolls []ActivePollResponse
	for _, poll := range polls {
		options, err := h.cfg.Queries.GetOptionsByPollID(r.Context(), poll.Pollid)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
			return
		}
		voteCount, err := h.cfg.Queries.GetTotalVotesByPollID(r.Context(), poll.Pollid)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
			return
		}
		commentCount, err := h.cfg.Queries.GetTotalComments(r.Context(), poll.Pollid)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
			return
		}
		pollResponse := ActivePollResponse{
			ID:          poll.Pollid,
			Title:       poll.Title,
			Creator:     poll.Creatorfirstname + " " + poll.Creatorlastname.String,
			Description: poll.Description,
			Category:    poll.Category,
			DaysLeft:    int64(poll.Expiresat.Sub(time.Now()).Hours() / 24),
			Options:     options,
			Votes:       voteCount,
			Comments:    commentCount,
		}

		ResponsePolls = append(ResponsePolls, pollResponse)
	}
	respondWithJSON(w, http.StatusOK, ResponsePolls)
	return
}
