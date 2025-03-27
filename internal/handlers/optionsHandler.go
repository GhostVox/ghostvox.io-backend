package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/GhostVox/ghostvox.io-backend/internal/config"
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
