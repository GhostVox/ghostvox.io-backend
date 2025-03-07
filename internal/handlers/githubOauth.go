package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/GhostVox/ghostvox.io-backend/internal/auth"
	"github.com/GhostVox/ghostvox.io-backend/internal/config"
	"github.com/GhostVox/ghostvox.io-backend/internal/database"

	"golang.org/x/oauth2"
)

type githubHandler struct {
	cfg               *config.APIConfig
	githubOAuthConfig *oauth2.Config
}

func NewGithubHandler(cfg *config.APIConfig, githubOAuthConfig *oauth2.Config) *githubHandler {
	return &githubHandler{
		cfg:               cfg,
		githubOAuthConfig: githubOAuthConfig,
	}
}

func (gh *githubHandler) GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := auth.GenerateRandomState()

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		HttpOnly: true,
		Secure:   true,
		Value:    state,
		Path:     "/",
		MaxAge:   300,
	})

	url := gh.githubOAuthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GithubCallbackHandler handles the OAuth callback from GitHub
func (gh *githubHandler) GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve state from query params
	stateParam := r.URL.Query().Get("state")

	// Retrieve state from cookie
	cookie, err := r.Cookie("oauth_state")
	if err != nil || cookie.Value != stateParam {
		errMsg := url.QueryEscape("Invalid OAuth state")
		http.Redirect(w, r, gh.cfg.AccessOrigin+"/sign-in?error="+errMsg, http.StatusTemporaryRedirect)
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
		errMsg := url.QueryEscape("Code not found")
		http.Redirect(w, r, gh.cfg.AccessOrigin+"/sign-in?error="+errMsg, http.StatusTemporaryRedirect)
		return
	}

	// Exchange auth code for access token
	token, err := gh.githubOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		errMsg := url.QueryEscape("Failed to exchange token")
		http.Redirect(w, r, gh.cfg.AccessOrigin+"/sign-in?error="+errMsg, http.StatusTemporaryRedirect)
		return
	}

	// Create a client with the token
	client := gh.githubOAuthConfig.Client(context.Background(), token)

	// Fetch user info from GitHub
	// GitHub API requires a specific Accept header
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		errMsg := url.QueryEscape("Failed to create request")
		http.Redirect(w, r, gh.cfg.AccessOrigin+"/sign-in?error="+errMsg, http.StatusTemporaryRedirect)
		return
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		errMsg := url.QueryEscape("Failed to get user info")
		http.Redirect(w, r, gh.cfg.AccessOrigin+"/sign-in?error="+errMsg, http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	// GitHub user info structure
	type GithubUser struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	}

	// Decode JSON response
	var githubUser GithubUser
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		errMsg := url.QueryEscape("Failed to decode user info")
		http.Redirect(w, r, gh.cfg.AccessOrigin+"/sign-in?error="+errMsg, http.StatusTemporaryRedirect)
		return
	}

	// If GitHub doesn't provide email in the user profile, we need to fetch emails separately
	if githubUser.Email == "" {
		emails, err := fetchGithubEmails(client)
		if err != nil {
			errMsg := url.QueryEscape("Failed to fetch user emails")
			http.Redirect(w, r, gh.cfg.AccessOrigin+"/sign-in?error="+errMsg, http.StatusTemporaryRedirect)
			return
		}

		// Use the primary email
		for _, email := range emails {
			if email.Primary {
				githubUser.Email = email.Email
				break
			}
		}
	}

	// Handle name, which might be empty
	var firstName, lastName string
	if githubUser.Name != "" {
		nameParts := strings.Split(githubUser.Name, " ")
		firstName = nameParts[0]
		if len(nameParts) > 1 {
			lastName = nameParts[1]
		}
	} else {
		// Use login as first name if name is not provided
		firstName = githubUser.Login
	}

	// Create user object
	user := User{
		FirstName:  firstName,
		LastName:   lastName,
		Email:      githubUser.Email,
		Provider:   githubProvider,
		PictureURL: githubUser.AvatarURL,
	}

	// Check if email already exists with a different provider
	existingUser, err := gh.cfg.Queries.GetUserByEmail(r.Context(), user.Email)
	if err == nil && existingUser.Provider.String != githubProvider {
		errMsg := url.QueryEscape("Account already exists with a different provider")
		http.Redirect(w, r, gh.cfg.AccessOrigin+"/sign-in?error="+errMsg, http.StatusTemporaryRedirect)
		return
	}

	var userRecord database.User
	var refreshTokenString string

	// If user doesn't exist, create a new one
	if errors.Is(err, sql.ErrNoRows) {
		refreshToken, newUserRecord, err := addUserAndRefreshToken(r.Context(), gh.cfg.DB, gh.cfg.Queries, &user)
		if err != nil {
			errMsg := url.QueryEscape("Failed to create user account")
			http.Redirect(w, r, gh.cfg.AccessOrigin+"/sign-in?error="+errMsg, http.StatusTemporaryRedirect)
			return
		}
		userRecord = newUserRecord
		refreshTokenString = refreshToken
	} else {

		refreshToken, err := deleteAndReplaceRefreshToken(r.Context(), gh.cfg, existingUser.ID)
		if err != nil {
			errMsg := url.QueryEscape("Failed to delete and replace refresh token")
			http.Redirect(w, r, gh.cfg.AccessOrigin+"/sign-in?error="+errMsg, http.StatusTemporaryRedirect)
			return
		}
		userRecord = existingUser
		refreshTokenString = refreshToken
	}

	// Generate Access Token
	accessToken, err := auth.GenerateJWTAccessToken(userRecord.ID, userRecord.Role, userRecord.PictureUrl.String, gh.cfg.GhostvoxSecretKey, gh.cfg.AccessTokenExp)
	if err != nil {
		errMsg := url.QueryEscape("Failed to generate access token")
		http.Redirect(w, r, gh.cfg.AccessOrigin+"/sign-in?error="+errMsg, http.StatusTemporaryRedirect)
		return
	}

	// Set cookies
	SetCookiesHelper(w, http.StatusOK, refreshTokenString, accessToken, gh.cfg)

	// Redirect to dashboard
	http.Redirect(w, r, gh.cfg.AccessOrigin+"/dashboard", http.StatusTemporaryRedirect)
}

// GitHub email structure
type GithubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

// Fetch user's GitHub emails
func fetchGithubEmails(client *http.Client) ([]GithubEmail, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Try to unmarshal as an array first
	var emails []GithubEmail
	err = json.Unmarshal(body, &emails)
	if err != nil {
		// If that fails, try unmarshaling as a single object
		var singleEmail GithubEmail
		if err := json.Unmarshal(body, &singleEmail); err != nil {
			return nil, err
		}
		emails = []GithubEmail{singleEmail}
	}

	return emails, nil
}
