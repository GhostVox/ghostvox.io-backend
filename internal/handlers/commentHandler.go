package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/GhostVox/ghostvox.io-backend/internal/auth"
	"github.com/GhostVox/ghostvox.io-backend/internal/config"
	"github.com/GhostVox/ghostvox.io-backend/internal/database"
	trie "github.com/Ghostvox/trie_hard/go"
	"github.com/google/uuid"
)

type CommentHandler struct {
	cfg    *config.APIConfig
	filter *trie.Trie[string]
}

type CommentResponse struct {
	ID        string         `json:"ID"`
	UserID    string         `json:"UserID"`
	UserName  sql.NullString `json:"Username"`
	AvatarUrl sql.NullString `json:"AvatarUrl"`
	Content   string         `json:"Content"`
	CreatedAt string         `json:"CreatedAt"`
}

func NewCommentHandler(cfg *config.APIConfig, filter *trie.Trie[string]) *CommentHandler {
	return &CommentHandler{
		cfg:    cfg,
		filter: filter,
	}
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
	cleanContent := cleanComment(comment.Content, h.filter)

	commentID, err := h.cfg.Queries.CreateComment(r.Context(), database.CreateCommentParams{
		UserID:  userUUID,
		PollID:  pollUUID,
		Content: cleanContent,
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
		return
	}
	commentUUID, err := uuid.Parse(commentID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "commentId", "Invalid CommentID", err)
		return

	}
	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "session", "Invalid session", err)
		return
	}
	if claims.Role == "admin" {
		err = h.cfg.Queries.AdminDeleteComment(r.Context(), commentUUID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "database", "Failed to delete comment", err)
			return
		}
		respondWithJSON(w, http.StatusNoContent, nil)
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

// Helper to clean comments using the trie filter
func cleanComment(content string, filter *trie.Trie[string]) string {
	contentSplit := strings.Fields(content)
	for i, word := range contentSplit {
		lower := strings.ToLower(word)
		cleanedWord := strings.TrimFunc(lower, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})

		if _, found := filter.Get(&cleanedWord); found {
			// Preserve prefix/suffix punctuation
			prefix := lower[:strings.Index(lower, cleanedWord)]
			suffix := lower[strings.Index(lower, cleanedWord)+len(cleanedWord):]
			contentSplit[i] = prefix + "****" + suffix
		}
	}
	return strings.Join(contentSplit, " ")
}
