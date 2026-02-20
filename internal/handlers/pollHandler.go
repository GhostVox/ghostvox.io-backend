package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/GhostVox/ghostvox.io-backend/internal/auth"
	"github.com/GhostVox/ghostvox.io-backend/internal/config"
	"github.com/GhostVox/ghostvox.io-backend/internal/database"
	t "github.com/Ghostvox/trie_hard/go"
	"github.com/google/uuid"
)

type poll struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Category    string         `json:"category"`
	ExpiresAt   string         `json:"expiresAt"`
	Status      string         `json:"status"`
	Options     []CreateOption `json:"options"`
}

type PollResponse struct {
	ID          uuid.UUID     `json:"id"`
	Title       string        `json:"title"`
	Creator     string        `json:"creator"`
	Description string        `json:"description"`
	Status      string        `json:"status"`
	Category    string        `json:"category"`
	DaysLeft    int64         `json:"daysLeft"`
	Options     []Option      `json:"options"`
	Votes       int64         `json:"votes"`
	Comments    int64         `json:"comments"`
	EndedAt     time.Time     `json:"endedAt"`
	Winner      string        `json:"winner"`
	UserVote    uuid.NullUUID `json:"userVote,omitempty"`
}

type pollHandler struct {
	cfg    *config.APIConfig
	filter *t.Trie[string]
}

func NewPollHandler(cfg *config.APIConfig, filter *t.Trie[string]) *pollHandler {
	return &pollHandler{
		cfg:    cfg,
		filter: filter,
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
}

func (h *pollHandler) GetPollByID(w http.ResponseWriter, r *http.Request, claims *auth.CustomClaims) {
	id := r.PathValue("pollId")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "pollId", "Missing poll ID", nil)
		return
	}
	pollID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "pollId", "Invalid poll ID", err)
		return
	}

	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "userUUID", "Invalid user UUID", err)
		return
	}

	poll, err := h.cfg.Queries.GetPollByID(r.Context(), database.GetPollByIDParams{
		ID:     pollID,
		UserID: userUUID,
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

	pollResponse, err := h.mapToPollResponse(
		poll.Pollid,
		poll.Title,
		poll.Description,
		poll.Category,
		string(poll.Status),
		poll.Creatorfirstname.String,
		poll.Creatorlastname.String,
		poll.Expiresat,
		poll.Votes,
		poll.Comments,
		poll.Options,
		poll.Uservote,
	)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "options", "Invalid options", err)
		return
	}

	respondWithJSON(w, http.StatusOK, pollResponse)
}

func (h *pollHandler) CreatePoll(w http.ResponseWriter, r *http.Request, claims *auth.CustomClaims) {

	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "accessToken", "Invalid access token", err)
		return
	}

	newPoll := poll{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&newPoll); err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
		return
	}

	// Check that title and description are present
	titlePresent, descriptionPresent := CheckPollTitleAndDescription(newPoll.Title, newPoll.Description)
	if !titlePresent {
		respondWithError(w, http.StatusBadRequest, "title", "Title is required", nil)
		return
	}
	if !descriptionPresent {
		newPoll.Description = "No description provided."
	}

	// Check for profanity in title, description, and options
	if !checkInputClean(newPoll.Title, h.filter, w) {
		return
	}
	if !checkInputClean(newPoll.Description, h.filter, w) {
		return
	}
	for _, option := range newPoll.Options {
		if !checkInputClean(option.Name, h.filter, w) {
			return
		}
	}

	err = CreatePollWithOptions(r.Context(), h.cfg, newPoll, userUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, struct {
		msg string
	}{msg: http.StatusText(http.StatusOK)})
}

func (h *pollHandler) UpdatePoll(w http.ResponseWriter, r *http.Request, claims *auth.CustomClaims) {
	pollId := r.PathValue("pollId")

	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "accessToken", "Invalid access token", err)
		return
	}

	newPoll := poll{}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&newPoll)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
		return
	}

	// Check that title and description are present
	titlePresent, descriptionPresent := CheckPollTitleAndDescription(newPoll.Title, newPoll.Description)
	if !titlePresent {
		respondWithError(w, http.StatusBadRequest, "title", "Title is required", nil)
		return
	}
	if !descriptionPresent {
		newPoll.Description = "No description provided."
	}
	// Validate inputs for profanity
	if !checkInputClean(newPoll.Title, h.filter, w) {
		return
	}
	if !checkInputClean(newPoll.Description, h.filter, w) {
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
		UserID:      userUUID,
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
}

func (h *pollHandler) DeletePoll(w http.ResponseWriter, r *http.Request, claims *auth.CustomClaims) {
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

	if claims.Role != "admin" {
		respondWithError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden), "Forbidden", nil)
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
}

func (h *pollHandler) GetAllActivePolls(w http.ResponseWriter, r *http.Request, claims *auth.CustomClaims) {
	limit, offset, err := getLimitAndOffset(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid limit or offset", err)
		return
	}

	category := r.URL.Query().Get("category")
	if category == "" {
		category = "%%"
	}

	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid user UUID", err)
		return
	}

	polls, err := h.cfg.Queries.GetAllPollsByStatusList(r.Context(), database.GetAllPollsByStatusListParams{
		Limit:    int32(limit),
		Offset:   int32(offset),
		Status:   database.PollStatus(database.PollStatusActive),
		Category: category,
		UserID:   userUUID,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithJSON(w, http.StatusOK, []PollResponse{})
			return
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
		return
	}
	pollsResp := make([]PollResponse, len(polls))
	for i, poll := range polls {

		p, err := h.mapToPollResponse(
			poll.Pollid,
			poll.Title,
			poll.Description,
			poll.Category,
			string(poll.Status),
			poll.Creatorfirstname,
			poll.Creatorlastname.String,
			poll.Expiresat,
			poll.Votes,
			poll.Comments,
			poll.Options,
			poll.Uservote,
		)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "options", "Invalid options", err)
			return
		}
		pollsResp[i] = p
	}

	respondWithJSON(w, http.StatusOK, pollsResp)

}

func (h *pollHandler) GetAllFinishedPolls(w http.ResponseWriter, r *http.Request, claims *auth.CustomClaims) {

	limit, offset, err := getLimitAndOffset(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid limit or offset", err)
		return
	}

	category := r.URL.Query().Get("category")
	if category == "" {
		category = "%%"
	}

	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid user UUID", err)
		return
	}

	polls, err := h.cfg.Queries.GetAllPollsByStatusList(r.Context(), database.GetAllPollsByStatusListParams{
		Limit:    int32(limit),
		Offset:   int32(offset),
		Status:   database.PollStatus(database.PollStatusArchived),
		Category: category,
		UserID:   userUUID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithJSON(w, http.StatusOK, []PollResponse{})
			return
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
		return
	}

	pollsResp := make([]PollResponse, len(polls))
	for i, poll := range polls {
		p, err := h.mapToPollResponse(
			poll.Pollid,
			poll.Title,
			poll.Description,
			poll.Category,
			string(poll.Status),
			poll.Creatorfirstname,
			poll.Creatorlastname.String,
			poll.Expiresat,
			poll.Votes,
			poll.Comments,
			poll.Options,
			poll.Uservote,
		)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "options", "Invalid options", err)
			return
		}
		pollsResp[i] = p
	}

	respondWithJSON(w, http.StatusOK, pollsResp)

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

	pollsResp := make([]PollResponse, len(userPolls))
	for i, poll := range userPolls {
		p, err := h.mapToPollResponse(
			poll.Pollid,
			poll.Title,
			poll.Description,
			poll.Category,
			string(poll.Status),
			poll.Creatorfirstname,
			poll.Creatorlastname.String,
			poll.Expiresat,
			poll.Votes,
			poll.Comments,
			poll.Options,
			poll.Uservote,
		)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "options", "Invalid options", err)
			return
		}
		pollsResp[i] = p
	}
	respondWithJSON(w, http.StatusOK, pollsResp)

}

func (h *pollHandler) GetRecentPolls(w http.ResponseWriter, r *http.Request, claims *auth.CustomClaims) {

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid user UUID", err)
		return
	}

	polls, err := h.cfg.Queries.GetRecentPolls(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithJSON(w, http.StatusOK, []PollResponse{})
			return
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Failed to retrieve recent polls", err)
		return
	}

	pollsResp := make([]PollResponse, len(polls))
	for i, poll := range polls {
		p, err := h.mapToPollResponse(
			poll.Pollid,
			poll.Title,
			poll.Description,
			poll.Category,
			string(poll.Status),
			poll.Creatorfirstname,
			poll.Creatorlastname.String,
			poll.Expiresat,
			poll.Votes,
			poll.Comments,
			poll.Options,
			poll.Uservote,
		)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "options", "Invalid options", err)
			return
		}

		pollsResp[i] = p
	}

	respondWithJSON(w, http.StatusOK, pollsResp)
}

// Create a helper to centralize the conversion logic
func (h *pollHandler) mapToPollResponse(
	pollID uuid.UUID,
	title, description, category, status string,
	creatorFirst, creatorLast string,
	expiresAt time.Time,
	votes, comments int64,
	optionsJSON []byte,
	userVote uuid.NullUUID,
) (PollResponse, error) {
	var options []Option
	if err := json.Unmarshal(optionsJSON, &options); err != nil {
		return PollResponse{}, err
	}

	return PollResponse{
		ID:          pollID,
		Title:       title,
		Creator:     creatorFirst + " " + creatorLast,
		Description: description,
		Status:      status,
		Category:    category,
		Options:     options,
		DaysLeft:    int64(time.Until(expiresAt).Hours() / 24),
		Votes:       votes,
		Comments:    comments,
		EndedAt:     expiresAt,
		Winner:      getWinner(options),
		UserVote:    userVote,
	}, nil
}

// Helper function to check for profanity in input
func checkInputClean(input string, filter *t.Trie[string], w http.ResponseWriter) bool {
	words := strings.Fields(input)
	for _, word := range words {
		cleanedWord := strings.TrimFunc(strings.ToLower(word), func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})

		if _, ok := filter.Get(&cleanedWord); ok {
			respondWithError(w, http.StatusBadRequest, "profanity", "Input contains profanity", nil)
			return false
		}
	}
	return true
}

// Helper to check if title and description are present
func CheckPollTitleAndDescription(title string, description string) (bool, bool) {
	titlePresent := true
	descriptionPresent := true
	if title == "" {
		titlePresent = false
	}
	if description == "" {
		descriptionPresent = false
	}
	return titlePresent, descriptionPresent
}
