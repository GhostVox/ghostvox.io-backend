package handlers

import (
	"context"
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
	ExpiresAt   string         `json:"expiresAt"`
	Status      string         `json:"status"`
	Options     []CreateOption `json:"options"`
}

type PollResponse struct {
	ID          uuid.UUID         `json:"id"`
	Title       string            `json:"title"`
	Creator     string            `json:"creator"`
	Description string            `json:"description"`
	Status      string            `json:"status"`
	Category    string            `json:"category"`
	DaysLeft    int64             `json:"daysLeft"`
	Options     []database.Option `json:"options"`
	Votes       int64             `json:"votes"`
	Comments    int64             `json:"comments"`
	EndedAt     time.Time         `json:"endedAt"`
	Winner      string            `json:"winner"`
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

func (h *pollHandler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	newPoll := poll{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&newPoll); err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
		return
	}

	err := CreatePollWithOptions(r.Context(), h.cfg, newPoll)
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
	exp, err := strconv.Atoi(newPoll.ExpiresAt)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid expiresAt", "Invalid expiresAt", err)
		return
	}
	expiresAt := time.Now().Add(time.Duration(exp))
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

// Helper function to process poll options, votes, and comments
func (h *pollHandler) processPollData(ctx context.Context, pollIDs []uuid.UUID) (
	map[uuid.UUID][]database.Option,
	map[uuid.UUID]int64,
	map[uuid.UUID]int64,
	error) {

	// Get all options, votes, and comments in batch operations
	allOptions, err := h.cfg.Queries.GetOptionsByPollIDs(ctx, pollIDs)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil, fmt.Errorf("failed to retrieve poll options: %w", err)
	}

	allVoteCounts, err := h.cfg.Queries.GetTotalVotesByPollIDs(ctx, pollIDs)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil, fmt.Errorf("failed to retrieve vote counts: %w", err)
	}

	allCommentCounts, err := h.cfg.Queries.GetTotalCommentsByPollIDs(ctx, pollIDs)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil, fmt.Errorf("failed to retrieve comment counts: %w", err)
	}

	// Organize data by poll ID for easy lookup
	optionsByPollID := make(map[uuid.UUID][]database.Option)
	for _, option := range allOptions {
		optionsByPollID[option.PollID] = append(optionsByPollID[option.PollID], option)
	}

	votesByPollID := make(map[uuid.UUID]int64)
	for _, vote := range allVoteCounts {
		votesByPollID[vote.PollID] = vote.Count
	}

	commentsByPollID := make(map[uuid.UUID]int64)
	for _, comment := range allCommentCounts {
		commentsByPollID[comment.PollID] = comment.Count
	}

	return optionsByPollID, votesByPollID, commentsByPollID, nil
}

// Process status-based polls
func (h *pollHandler) processStatusPolls(ctx context.Context, polls []database.GetAllPollsByStatusListRow, w http.ResponseWriter) {
	if len(polls) == 0 {
		respondWithJSON(w, http.StatusOK, []PollResponse{})
		return
	}

	// Extract all poll IDs
	pollIDs := make([]uuid.UUID, len(polls))
	for i, poll := range polls {
		pollIDs[i] = poll.Pollid
	}

	// Get processed poll data
	optionsByPollID, votesByPollID, commentsByPollID, err := h.processPollData(ctx, pollIDs)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
		return
	}

	// Build response
	responsePolls := make([]PollResponse, 0, len(polls))
	for _, poll := range polls {
		pollID := poll.Pollid
		options := optionsByPollID[pollID]

		pollResponse := PollResponse{
			ID:          pollID,
			Title:       poll.Title,
			Creator:     poll.Creatorfirstname + " " + poll.Creatorlastname.String,
			Status:      string(poll.Status),
			Description: poll.Description,
			Category:    poll.Category,
			Options:     options,
			DaysLeft:    int64(poll.Expiresat.Sub(time.Now()).Hours() / 24),
			Votes:       votesByPollID[pollID],
			Comments:    commentsByPollID[pollID],
			EndedAt:     poll.Expiresat,
		}

		if poll.Status == database.PollStatusArchived && len(options) > 0 {
			pollResponse.Winner = getWinner(options)
		}

		responsePolls = append(responsePolls, pollResponse)
	}

	respondWithJSON(w, http.StatusOK, responsePolls)
}

// Process user-based polls
func (h *pollHandler) processUserPolls(ctx context.Context, polls []database.GetPollsByUserRow, w http.ResponseWriter) {
	if len(polls) == 0 {
		respondWithJSON(w, http.StatusOK, []PollResponse{})
		return
	}

	// Extract all poll IDs
	pollIDs := make([]uuid.UUID, len(polls))
	for i, poll := range polls {
		pollIDs[i] = poll.Pollid
	}

	// Get processed poll data
	optionsByPollID, votesByPollID, commentsByPollID, err := h.processPollData(ctx, pollIDs)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
		return
	}

	// Build response
	responsePolls := make([]PollResponse, 0, len(polls))
	for _, poll := range polls {
		pollID := poll.Pollid
		options := optionsByPollID[pollID]

		pollResponse := PollResponse{
			ID:          pollID,
			Title:       poll.Title,
			Creator:     poll.Creatorfirstname + " " + poll.Creatorlastname.String,
			Status:      string(poll.Status),
			Description: poll.Description,
			Category:    poll.Category,
			Options:     options,
			DaysLeft:    int64(poll.Expiresat.Sub(time.Now()).Hours() / 24),
			Votes:       votesByPollID[pollID],
			Comments:    commentsByPollID[pollID],
			EndedAt:     poll.Expiresat,
		}

		if poll.Status == database.PollStatusArchived && len(options) > 0 {
			pollResponse.Winner = getWinner(options)
		}

		responsePolls = append(responsePolls, pollResponse)
	}

	respondWithJSON(w, http.StatusOK, responsePolls)
}

func (h *pollHandler) GetAllActivePolls(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := getLimitAndOffset(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid limit or offset", err)
		return
	}

	category := r.URL.Query().Get("category")
	if category == "" {
		category = "%%"
	}

	polls, err := h.cfg.Queries.GetAllPollsByStatusList(r.Context(), database.GetAllPollsByStatusListParams{
		Limit:    int32(limit),
		Offset:   int32(offset),
		Status:   database.PollStatus(database.PollStatusActive),
		Category: category,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithJSON(w, http.StatusOK, []PollResponse{})
			return
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
		return
	}

	h.processStatusPolls(r.Context(), polls, w)
}

func (h *pollHandler) GetAllFinishedPolls(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := getLimitAndOffset(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid limit or offset", err)
		return
	}

	category := r.URL.Query().Get("category")
	if category == "" {
		category = "%%"
	}

	polls, err := h.cfg.Queries.GetAllPollsByStatusList(r.Context(), database.GetAllPollsByStatusListParams{
		Limit:    int32(limit),
		Offset:   int32(offset),
		Status:   database.PollStatus(database.PollStatusArchived),
		Category: category,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithJSON(w, http.StatusOK, []PollResponse{})
			return
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
		return
	}

	h.processStatusPolls(r.Context(), polls, w)
}

func (h *pollHandler) GetUsersPolls(w http.ResponseWriter, r *http.Request) {
	userIDString := r.PathValue("userId")
	userId, err := uuid.Parse(userIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid user ID", err)
		return
	}

	limit, offset, err := getLimitAndOffset(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid limit or offset", err)
		return
	}

	category := r.URL.Query().Get("category")
	if category == "" {
		category = "%%"
	}

	userPolls, err := h.cfg.Queries.GetPollsByUser(r.Context(), database.GetPollsByUserParams{
		UserID:   userId,
		Limit:    int32(limit),
		Offset:   int32(offset),
		Category: category,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithJSON(w, http.StatusOK, []PollResponse{})
			return
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Failed to retrieve user polls", err)
		return
	}

	h.processUserPolls(r.Context(), userPolls, w)
}
