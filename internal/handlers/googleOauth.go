package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/GhostVox/ghostvox.io-backend/internal/auth"
	"github.com/GhostVox/ghostvox.io-backend/internal/config"
	"github.com/GhostVox/ghostvox.io-backend/internal/database"

	"golang.org/x/oauth2"
)

const (
	googleProvider = "google"
	githubProvider = "github"
)

type googleHandler struct {
	cfg               *config.APIConfig
	googleOAuthConfig *oauth2.Config
}

// NewGoogleHandler
func NewGoogleHandler(cfg *config.APIConfig, googleOAuthConfig *oauth2.Config) *googleHandler {
	return &googleHandler{
		cfg:               cfg,
		googleOAuthConfig: googleOAuthConfig,
	}
}

// Generate Google Login URL
func (gh *googleHandler) GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Generate a random state token
	state := auth.GenerateRandomState()

	// Store state token in a secure cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
		Secure:   true, // Set to true in production with HTTPS
		Path:     "/",
		MaxAge:   300, // 5 minutes expiration
	})

	// Redirect user to Google's OAuth URL with the state
	url := gh.googleOAuthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Handle Google OAuth Callback
func (gh *googleHandler) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve state from query params
	stateParam := r.URL.Query().Get("state")

	// Retrieve state from cookie
	cookie, err := r.Cookie("oauth_state")
	if err != nil || cookie.Value != stateParam {
		http.Error(w, "Invalid OAuth state", http.StatusUnauthorized)
		return
	}

	// Clear state cookie to prevent reuse
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   -1, // Deletes the cookie
	})
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	// Exchange auth code for access token
	token, err := gh.googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch user info from Google
	client := gh.googleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Decode JSON response
	var oAuthUser config.OAuthUser
	if err := json.NewDecoder(resp.Body).Decode(&oAuthUser); err != nil {
		http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Split Name into First & Last Name
	nameParts := strings.Split(oAuthUser.Name, " ")
	firstName := nameParts[0]
	lastName := ""
	if len(nameParts) > 1 {
		lastName = nameParts[1]
	}

	user := User{
		FirstName:  firstName,
		LastName:   lastName,
		Email:      oAuthUser.Email,
		Provider:   googleProvider,
		PictureURL: oAuthUser.PictureURL,
	}

	// Check if email already exists (prevents duplicate accounts)
	existingUser, err := gh.cfg.Queries.GetUserByEmail(r.Context(), user.Email)
	if err == nil && existingUser.Provider.String != googleProvider {
		http.Error(w, "Account already exists with a different provider", http.StatusConflict)
		return
	}

	var userRecord database.User
	var refreshTokenString string
	// If user doesn't exist, create a new one
	if errors.Is(err, sql.ErrNoRows) {
		refreshToken, newUserRecord, err := addUserAndRefreshToken(r.Context(), gh.cfg.DB, gh.cfg.Queries, user)
		if err != nil {
			http.Error(w, "Failed to add user and refresh token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		userRecord = newUserRecord
		refreshTokenString = refreshToken
	} else {
		refreshToken, err := auth.GenerateRefreshToken()
		if err != nil {
			http.Error(w, "Failed to generate refresh token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		userRecord = existingUser
		refreshTokenString = refreshToken
	}
	// Generate Access Token
	accessToken, err := auth.GenerateJWTAccessToken(userRecord.ID, userRecord.Role, userRecord.PictureUrl.String, gh.cfg.GhostvoxSecretKey, gh.cfg.AccessTokenExp)
	if err != nil {
		http.Error(w, "Failed to generate access token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set cookies
	SetCookiesHelper(w, refreshTokenString, accessToken, gh.cfg)
}
