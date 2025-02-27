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
	Value     string `json:"value"`
	PollID    string `json:"poll_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
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
		chooseError(w, http.StatusBadRequest, errors.New("option id is required"))
		return
	}
	optionUUID, err := uuid.Parse(id)
	if err != nil {
		chooseError(w, http.StatusBadRequest, errors.New("invalid option id"))
		return
	}
	optionRecord, err := oh.cfg.DB.GetOptionByID(r.Context(), optionUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusNotFound, errors.New("option not found"))
			return
		}
		chooseError(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, http.StatusOK, optionRecord)

	return
}

func (oh *optionHandler) GetOptionsByPollID(w http.ResponseWriter, r *http.Request) {
	pollID := r.PathValue("pollId")
	if pollID == "" {
		chooseError(w, http.StatusBadRequest, errors.New("poll id is required"))
		return
	}
	pollUUID, err := uuid.Parse(pollID)
	if pollUUID == uuid.Nil {
		chooseError(w, http.StatusBadRequest, errors.New("invalid poll id"))
		return
	}
	options, err := oh.cfg.DB.GetOptionsByPollID(r.Context(), pollUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusNotFound, errors.New("options not found"))
			return
		}
		chooseError(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, http.StatusOK, options)

	return
}

func (oh *optionHandler) CreateOptions(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("pollId")
	if pollId == "" {
		chooseError(w, http.StatusBadRequest, errors.New("poll id is required"))
		return
	}
	pollUUID, err := uuid.Parse(pollId)
	if pollUUID == uuid.Nil {
		chooseError(w, http.StatusBadRequest, errors.New("invalid poll id"))
		return
	}

	var pollOptions OptionsRequest
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&pollOptions)
	if err != nil {
		chooseError(w, http.StatusBadRequest, err)
		return
	}
	OptionRecords := []database.CreateOptionRow{}
	for _, option := range pollOptions.Options {
		optionRecord, err := oh.cfg.DB.CreateOption(r.Context(), database.CreateOptionParams{
			Name:   option.Name,
			Value:  option.Value,
			PollID: pollUUID,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				chooseError(w, http.StatusConflict, err)
				return
			}
			chooseError(w, http.StatusInternalServerError, err)
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
		chooseError(w, http.StatusBadRequest, errors.New("poll id is required"))
		return
	}
	pollUUID, err := uuid.Parse(pollId)
	if pollUUID == uuid.Nil {
		chooseError(w, http.StatusBadRequest, errors.New("invalid poll id"))
		return
	}

	var option Option
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&option)
	if err != nil {
		chooseError(w, http.StatusBadRequest, err)
		return
	}

	optionUUID, err := uuid.Parse(option.ID)
	if err != nil {
		chooseError(w, http.StatusBadRequest, err)
		return
	}

	optionRecord, err := oh.cfg.DB.UpdateOption(r.Context(), database.UpdateOptionParams{
		ID:    optionUUID,
		Name:  option.Name,
		Value: option.Value,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusNotFound, err)
			return
		}
		chooseError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, optionRecord)
	return
}

func (oh *optionHandler) DeleteOption(w http.ResponseWriter, r *http.Request) {
	optionId := r.PathValue("optionId")
	if optionId == "" {
		chooseError(w, http.StatusBadRequest, errors.New("option id is required"))
		return
	}
	optionUUID, err := uuid.Parse(optionId)
	if optionUUID == uuid.Nil {
		chooseError(w, http.StatusBadRequest, errors.New("invalid option id"))
		return
	}

	err = oh.cfg.DB.DeleteOption(r.Context(), optionUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusNotFound, err)
			return
		}
		chooseError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
