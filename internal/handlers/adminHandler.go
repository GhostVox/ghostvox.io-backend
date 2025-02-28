package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/GhostVox/ghostvox.io-backend/internal/config"
	"github.com/google/uuid"
)

type AdminHandler struct {
	cfg *config.APIConfig
}

func NewAdminHandler(cfg *config.APIConfig) *AdminHandler {
	return &AdminHandler{
		cfg: cfg,
	}
}

func (h *AdminHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	users, err := h.cfg.Queries.GetUsers(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusNotFound, err)
			return
		}
		chooseError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, users)
	return
}

func (h *AdminHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("userId")
	userUUID, err := uuid.Parse(id)
	if err != nil {
		chooseError(w, http.StatusBadRequest, err)
		return
	}
	user, err := h.cfg.Queries.GetUserById(r.Context(), userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusNotFound, err)
			return
		}
		chooseError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
	return
}
