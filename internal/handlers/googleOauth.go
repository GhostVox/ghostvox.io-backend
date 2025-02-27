package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/GhostVox/ghostvox.io-backend/internal/auth"
	"github.com/GhostVox/ghostvox.io-backend/internal/config"
	"github.com/GhostVox/ghostvox.io-backend/internal/database"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	googleProvider = "google"
	githubProvider = "github"
)

type OAuthUser struct {
	Email        string `json:"email,omitempty"`
	Name         string `json:"name,omitempty"`
	Password     string `json:"password,omitempty"`
	Provider     string `json:"provider,omitempty"`
	ProviderID   string `json:"id,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Role         string `json:"role,omitempty"`
	PictureURL   string `json:"picture,omitempty"`
}

// OAuth2 configuration
var googleOAuthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:3000/auth/callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

type googleHandler struct {
	cfg *config.APIConfig
}

// GoogleHandler
func NewGoogleHandler(cfg *config.APIConfig) *googleHandler {
	return &googleHandler{
		cfg: cfg,
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
	url := googleOAuthConfig.AuthCodeURL(state)
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
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch user info from Google
	client := googleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Decode JSON response
	var user OAuthUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Split Name into First & Last Name
	nameParts := strings.Split(user.Name, " ")
	firstName := nameParts[0]
	lastName := ""
	if len(nameParts) > 1 {
		lastName = nameParts[1]
	}

	// Check if email already exists (prevents duplicate accounts)
	_, err = gh.cfg.DB.GetUserByEmail(r.Context(), user.Email)
	if err == nil {
		http.Error(w, "Account already exists with a different provider", http.StatusConflict)
		return
	}

	// Check if user exists by provider & provider_id
	userRecord, err := gh.cfg.DB.GetUserByProviderAndProviderId(r.Context(), database.GetUserByProviderAndProviderIdParams{
		Provider:   NullStringHelper(googleProvider),
		ProviderID: NullStringHelper(user.ProviderID),
	})

	// If user doesn't exist, create a new one
	if errors.Is(err, sql.ErrNoRows) {
		newUserRecord, err := gh.cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
			Email:      user.Email,
			FirstName:  firstName,
			LastName:   NullStringHelper(lastName),
			Provider:   NullStringHelper(googleProvider),
			ProviderID: NullStringHelper(user.ProviderID),
			PictureUrl: NullStringHelper(user.PictureURL),
		})
		if err != nil {
			http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		userRecord = newUserRecord
	}

	// Generate Refresh Token
	refreshToken, err := AddRefreshToken(r.Context(), userRecord.ID, gh.cfg.DB)
	if err != nil {
		http.Error(w, "Failed to add refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate Access Token
	accessToken, err := auth.GenerateJWTAccessToken(userRecord.ID, userRecord.Role, userRecord.PictureUrl.String, gh.cfg.GhostvoxSecretKey, gh.cfg.AccessTokenExp)
	if err != nil {
		http.Error(w, "Failed to generate access token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set cookies
	SetCookiesHelper(w, refreshToken, accessToken, gh.cfg)
}
