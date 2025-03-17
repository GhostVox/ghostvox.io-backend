package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/GhostVox/ghostvox.io-backend/internal/config"
	"github.com/GhostVox/ghostvox.io-backend/internal/database"
	"github.com/google/uuid"
)

type Option struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PollID    string `json:"poll_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
type CreateOption struct {
	Name string `json:"name"`
}
type OptionsRequest struct {
	Options []Option `json:"options"`
}

type optionHandler struct {
	cfg *config.APIConfig
}

func NewOptionHandler(cfg *config.APIConfig) *optionHandler {
	return &optionHandler{cfg: cfg}
}

func (oh *optionHandler) GetOptionByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("optionId")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Missing option id pathvalue", errors.New("Missing option id pathvalue"))
		return
	}
	optionUUID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid option id", errors.New("Invalid option id"))
		return
	}
	optionRecord, err := oh.cfg.Queries.GetOptionByID(r.Context(), optionUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), "Option not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
		return
	}
	respondWithJSON(w, http.StatusOK, optionRecord)

	return
}

func (oh *optionHandler) GetOptionsByPollID(w http.ResponseWriter, r *http.Request) {
	pollID := r.PathValue("pollId")
	if pollID == "" {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Poll id is required", errors.New("Poll id is required"))
		return
	}
	pollUUID, err := uuid.Parse(pollID)
	if pollUUID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid poll id", errors.New("Invalid poll id"))
		return
	}
	options, err := oh.cfg.Queries.GetOptionsByPollID(r.Context(), pollUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), "Options not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
		return
	}
	respondWithJSON(w, http.StatusOK, options)

	return
}

func (oh *optionHandler) CreateOptions(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("pollId")
	if pollId == "" {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Poll id is required", errors.New("Poll id is required"))
		return
	}
	pollUUID, err := uuid.Parse(pollId)
	if pollUUID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid poll id", errors.New("Invalid poll id"))
		return
	}

	var pollOptions OptionsRequest
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&pollOptions)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid request body", err)
		return
	}
	OptionRecords := []database.CreateOptionRow{}
	for _, option := range pollOptions.Options {
		optionRecord, err := oh.cfg.Queries.CreateOption(r.Context(), database.CreateOptionParams{
			Name:   option.Name,
			PollID: pollUUID,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondWithError(w, http.StatusConflict, http.StatusText(http.StatusConflict), "Option already exists", err)
				return
			}
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
			return
		}
		OptionRecords = append(OptionRecords, optionRecord)
	}
	respondWithJSON(w, http.StatusCreated, OptionRecords)
	return
}

func (oh *optionHandler) UpdateOption(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("pollId")
	if pollId == "" {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Poll id is required", errors.New("Poll id is required"))
		return
	}
	pollUUID, err := uuid.Parse(pollId)
	if pollUUID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid poll id", errors.New("Invalid poll id"))
		return
	}

	var option Option
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&option)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid request body", err)
		return
	}

	optionUUID, err := uuid.Parse(option.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid option id", err)
		return
	}

	optionRecord, err := oh.cfg.Queries.UpdateOption(r.Context(), database.UpdateOptionParams{
		ID:   optionUUID,
		Name: option.Name,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), "Option not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
		return
	}

	respondWithJSON(w, http.StatusOK, optionRecord)
	return
}

func (oh *optionHandler) DeleteOption(w http.ResponseWriter, r *http.Request) {
	optionId := r.PathValue("optionId")
	if optionId == "" {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Option id is required", errors.New("Option id is required"))
		return
	}
	optionUUID, err := uuid.Parse(optionId)
	if optionUUID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid option id", errors.New("Invalid option id"))
		return
	}

	err = oh.cfg.Queries.DeleteOption(r.Context(), optionUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), "Option not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
