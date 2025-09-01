package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/GhostVox/ghostvox.io-backend/internal/auth"
	"github.com/GhostVox/ghostvox.io-backend/internal/config"
	"github.com/GhostVox/ghostvox.io-backend/internal/database"
	"github.com/google/uuid"
)

type CommentHandler struct {
	cfg *config.APIConfig
}

type CommentResponse struct {
	ID        string         `json:"ID"`
	UserID    string         `json:"UserID"`
	UserName  sql.NullString `json:"Username"`
	AvatarUrl sql.NullString `json:"AvatarUrl"`
	Content   string         `json:"Content"`
	CreatedAt string         `json:"CreatedAt"`
}

func NewCommentHandler(cfg *config.APIConfig) *CommentHandler {
	return &CommentHandler{cfg: cfg}
}

func (h *CommentHandler) GetAllPollComments(w http.ResponseWriter, r *http.Request) {
	pollID := r.PathValue("pollId")
	if pollID == "" {
		respondWithError(w, http.StatusBadRequest, "pollId", "Poll ID is required", nil)
		return
	}

	pollUUID, err := uuid.Parse(pollID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "pollId", "Invalid Poll ID", err)
		return
	}

	comments, err := h.cfg.Queries.GetAllCommentsByPollID(r.Context(), pollUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithJSON(w, http.StatusOK, []CommentResponse{})
			return
		}
		respondWithError(w, http.StatusInternalServerError, "database", "Failed to retrieve comments", err)
		return
	}
	respondWithJSON(w, http.StatusOK, comments)

}

func (h *CommentHandler) CreatePollComment(w http.ResponseWriter, r *http.Request, claims *auth.CustomClaims) {
	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "session", "Invalid session", err)
		return
	}

	pollID := r.PathValue("pollId")
	if pollID == "" {
		respondWithError(w, http.StatusBadRequest, "pollId", "Poll ID is required", nil)
		return
	}

	pollUUID, err := uuid.Parse(pollID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "pollId", "Invalid Poll ID", err)
		return
	}

	var comment struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		respondWithError(w, http.StatusBadRequest, "content", "Invalid content", err)
		return
	}

	if comment.Content == "" {
		respondWithError(w, http.StatusBadRequest, "content", "Content is required", nil)
		return
	}

	commentID, err := h.cfg.Queries.CreateComment(r.Context(), database.CreateCommentParams{
		UserID:  userUUID,
		PollID:  pollUUID,
		Content: comment.Content,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "database", "Failed to create comment", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, CommentResponse{
		ID:        commentID.String(),
		UserID:    userUUID.String(),
		UserName:  NullStringHelper(claims.UserName),
		AvatarUrl: NullStringHelper(claims.PictureUrl),
		Content:   comment.Content,
		CreatedAt: time.Now().Format(time.RFC3339),
	})

}

func (h *CommentHandler) DeletePollComment(w http.ResponseWriter, r *http.Request, claims *auth.CustomClaims) {
	// To be implemented
	commentID := r.PathValue("commentId")
	if commentID == "" {

		respondWithError(w, http.StatusBadRequest, "commentId", "Comment ID is required", nil)
	}
	commentUUID, err := uuid.Parse(commentID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "commentId", "Invalid CommentID", err)
	}
	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "session", "Invalid session", err)
		return
	}

	err = h.cfg.Queries.DeleteComment(r.Context(), database.DeleteCommentParams{
		ID: commentUUID, UserID: userUUID})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "database", "Failed to delete comment", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
