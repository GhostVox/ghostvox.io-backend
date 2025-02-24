package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/GhostVox/ghostvox.io-backend/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"name"`
	Email     string    `json:"email"`
	LastName  string    `json:"last_name"`
	UserToken string    `json:"user_token"`
	Role      string    `json:"role"`
}

type UserHandler struct {
	db *database.Queries
}

func NewUserHandler(db *database.Queries) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.db.GetUsers(r.Context())
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

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("userId")
	if id == "" {
		chooseError(w, http.StatusBadRequest, errors.New("missing id"))
		return
	}
	UserUUID, err := uuid.Parse(id)
	if err != nil {
		chooseError(w, http.StatusBadRequest, errors.New("invalid id"))
		return
	}
	user, err := h.db.GetUserById(r.Context(), UserUUID)
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

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		chooseError(w, http.StatusInternalServerError, err)
		return
	}

	userRecord, err := h.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:     user.Email,
		LastName:  sql.NullString{String: user.LastName, Valid: user.LastName != ""},
		FirstName: user.FirstName,
		Role:      user.Role,
		UserToken: user.UserToken,
	})
	if err != nil {
		chooseError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, userRecord)
	return
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("userId")
	if id == "" {
		chooseError(w, http.StatusBadRequest, errors.New("missing id"))
		return
	}

	var user User
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		chooseError(w, http.StatusInternalServerError, err)
		return
	}
	UserUUID, err := uuid.Parse(id)
	if err != nil {
		chooseError(w, http.StatusBadRequest, errors.New("invalid id"))
		return
	}
	userRecord, err := h.db.GetUserById(r.Context(), UserUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusNotFound, err)
			return
		}
		chooseError(w, http.StatusInternalServerError, err)
		return
	}

	updatedUserRecord, err := h.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:        userRecord.ID,
		Email:     user.Email,
		LastName:  sql.NullString{String: user.LastName, Valid: user.LastName != ""},
		FirstName: user.FirstName,
		Role:      user.Role,
		UserToken: user.UserToken,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusNotFound, err)
			return
		}
		chooseError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, updatedUserRecord)
	return
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("userId")
	if id == "" {
		chooseError(w, http.StatusBadRequest, errors.New("missing id"))
		return
	}

	userUUID, err := uuid.Parse(id)
	if err != nil {
		chooseError(w, http.StatusBadRequest, err)
		return
	}

	err = h.db.DeleteUser(r.Context(), userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusNotFound, err)
			return
		}
		chooseError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
	return
}
