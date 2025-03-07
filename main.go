package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"database/sql"

	"github.com/GhostVox/ghostvox.io-backend/internal/config"
	"github.com/GhostVox/ghostvox.io-backend/internal/database"
	"github.com/GhostVox/ghostvox.io-backend/internal/handlers"
	mw "github.com/GhostVox/ghostvox.io-backend/internal/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

func main() {

	const (
		port = ":8080"
	)

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db_URL := os.Getenv("DB_URL")
	if db_URL == "" {
		log.Fatal("DB_URL must be set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	ghostvoxSecretKey := os.Getenv("GHOSTVOX_SECRET_KEY")
	if ghostvoxSecretKey == "" {
		log.Fatal("GHOSTVOX_SECRET_KEY must be set")
	}

	accesstokenExp := os.Getenv("ACCESS_TOKEN_EXPIRES")
	if accesstokenExp == "" {
		log.Fatal("ACCESS_TOKEN_EXPIRES must be set")
	}

	refreshTokenExp := os.Getenv("REFRESH_TOKEN_EXPIRES")
	if refreshTokenExp == "" {
		log.Fatal("REFRESH_TOKEN_EXPIRES must be set")
	}
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	if googleClientID == "" {
		log.Fatal("GOOGLE_CLIENT_ID must be set")
	}
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if googleClientSecret == "" {
		log.Fatal("GOOGLE_CLIENT_SECRET must be set")
	}
	googleRedirectURI := os.Getenv("GOOGLE_REDIRECT_URI")
	if googleRedirectURI == "" {
		log.Fatal("GOOGLE_REDIRECT_URI must be set")
	}
	githubClientID := os.Getenv("GITHUB_CLIENT_ID")
	if githubClientID == "" {
		log.Fatal("GITHUB_CLIENT_ID must be set")
	}
	githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	if githubClientSecret == "" {
		log.Fatal("GITHUB_CLIENT_SECRET must be set")
	}
	githubRedirectURI := os.Getenv("GITHUB_REDIRECT_URI")
	if githubRedirectURI == "" {
		log.Fatal("GITHUB_REDIRECT_URI must be set")
	}
	mode := os.Getenv("MODE")
	if mode == "" {
		log.Fatal("MODE must be set")
	}

	https := os.Getenv("USE_HTTPS")
	if https == "" {
		log.Fatal("HTTPS must be set")
	}
	accessOrigin := os.Getenv("ACCESS_ORIGIN")
	if accessOrigin == "" {
		log.Fatal("ACCESS_ORIGIN must be set")
	}

	accesstokenExpDur, err := time.ParseDuration(accesstokenExp)
	if err != nil {
		log.Fatal(err)
	}

	refreshTokenExpDur, err := time.ParseDuration(refreshTokenExp)
	if err != nil {
		log.Fatal(err)
	}

	//Connect to database
	db, err := sql.Open("postgres", db_URL)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error connecting to database: %v", err))
	}

	//Connect database to sqlc code to get a pointer to queries to build handlers
	dbConnection := database.New(db)
	if err != nil {
		log.Fatal(err)
	}
	//Configure API struct to pass around
	cfg := &config.APIConfig{
		DB:                db,
		Queries:           dbConnection,
		Platform:          platform,
		Port:              port,
		AccessTokenExp:    accesstokenExpDur,
		RefreshTokenExp:   refreshTokenExpDur,
		GhostvoxSecretKey: ghostvoxSecretKey,
		Mode:              mode,
		UseHTTPS:          https,
		AccessOrigin:      accessOrigin,
	}
	// OAuth2 configuration
	var googleOAuthConfig = &oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		RedirectURL:  googleRedirectURI,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	var githubOAuthConfig = &oauth2.Config{
		ClientID:     githubClientID,
		ClientSecret: githubClientSecret,
		RedirectURL:  githubRedirectURI,
		Scopes: []string{
			"user:email",
			"read:user",
		},
		Endpoint: github.Endpoint,
	}
	// Initialize handlers
	rootHandler := handlers.NewRootHandler(cfg)
	pollHandler := handlers.NewPollHandler(cfg.Queries)
	voteHandler := handlers.NewVoteHandler(cfg.Queries)
	optionHandler := handlers.NewOptionHandler(cfg)
	userHandler := handlers.NewUserHandler(cfg)
	authHandler := handlers.NewAuthHandler(cfg)
	googleHandler := handlers.NewGoogleHandler(cfg, googleOAuthConfig)
	githubHandler := handlers.NewGithubHandler(cfg, githubOAuthConfig)
	adminHandler := handlers.NewAdminHandler(cfg)

	mux := http.NewServeMux()

	wrappedMux := mw.CorsMiddleware(mux, accessOrigin)
	//  start attaching route handlers to cfg.mux
	// Redirects users to the root of the API and returns route endpoints for the API
	mux.HandleFunc("/api/v1/", mw.LoggingMiddleware(rootHandler.HandleRoot))

	// Polls route ✅
	mux.HandleFunc("GET /api/v1/polls", mw.LoggingMiddleware(pollHandler.GetAllPolls))

	mux.HandleFunc("GET /api/v1/polls/{pollId}", mw.LoggingMiddleware(pollHandler.GetPoll))

	mux.HandleFunc("PUT /api/v1/polls/{pollId}", mw.LoggingMiddleware(pollHandler.UpdatePoll))

	mux.HandleFunc("POST /api/v1/polls", mw.LoggingMiddleware(pollHandler.CreatePoll))

	mux.HandleFunc("DELETE /api/v1/polls/{pollId}", mw.LoggingMiddleware(pollHandler.DeletePoll))
	// End of poll routes
	// OAuth routes
	mux.HandleFunc("GET /api/v1/auth/google/login", mw.LoggingMiddleware(googleHandler.GoogleLoginHandler))
	mux.HandleFunc("GET /api/v1/auth/google/callback", mw.LoggingMiddleware(googleHandler.GoogleCallbackHandler))
	mux.HandleFunc("GET /api/v1/auth/github/login", mw.LoggingMiddleware(githubHandler.GithubLoginHandler))
	mux.HandleFunc("GET /api/v1/auth/github/callback", mw.LoggingMiddleware(githubHandler.GithubCallbackHandler))
	// end
	//Auth routes
	mux.HandleFunc("POST /api/v1/auth/login", mw.LoggingMiddleware(authHandler.Login))
	mux.HandleFunc("POST /api/v1/auth/register", mw.LoggingMiddleware(authHandler.Register))
	mux.HandleFunc("POST /api/v1/auth/logout", mw.LoggingMiddleware(authHandler.Logout))
	mux.HandleFunc("POST /api/v1/auth/refresh", mw.LoggingMiddleware(authHandler.Refresh))

	// Users Private route
	mux.HandleFunc("GET /api/v1/admin/users/{userId}", mw.AdminRole(cfg, mw.LoggingMiddleware(adminHandler.GetUser)).ServeHTTP)

	mux.HandleFunc("GET /api/v1/admin/users", mw.AdminRole(cfg, mw.LoggingMiddleware(adminHandler.GetAllUsers)).ServeHTTP)

	// User public route ✅

	mux.HandleFunc("PUT /api/v1/users/{userId}", mw.LoggingMiddleware(userHandler.UpdateUser))

	mux.HandleFunc("DELETE /api/v1/users/{userId}", mw.LoggingMiddleware(userHandler.DeleteUser))
	// End of users routes

	// Votes Routes

	mux.HandleFunc("POST /api/v1/polls/{pollId}/votes", mw.LoggingMiddleware(voteHandler.CreateVote))

	mux.HandleFunc("GET /api/v1/polls/{pollId}/votes", mw.LoggingMiddleware(voteHandler.GetVotesByPoll))

	mux.HandleFunc("DELETE /api/v1/votes/{voteId}", mw.LoggingMiddleware(voteHandler.DeleteVote))

	// Options Routes
	mux.HandleFunc("POST /api/v1/polls/{pollId}/options", mw.LoggingMiddleware(optionHandler.CreateOptions))

	mux.HandleFunc("GET /api/v1/polls/{pollId}/options/{optionId}", mw.LoggingMiddleware(optionHandler.GetOptionByID))

	mux.HandleFunc("GET /api/v1/polls/{pollId}/options", mw.LoggingMiddleware(optionHandler.GetOptionsByPollID))

	mux.HandleFunc("PUT /api/v1/polls/{pollId}/options/{optionId}", mw.LoggingMiddleware(optionHandler.UpdateOption))

	mux.HandleFunc("DELETE /api/v1/polls/{pollId}/options/{optionId}", mw.LoggingMiddleware(optionHandler.DeleteOption))
	// End of options routes

	//Redirect from root page to API documentation
	mux.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "http://localhost:8080/api/v1", http.StatusFound)
			return
		}

		http.NotFound(w, r)
	}))

	addr := port
	// Create a pointer to the http.Server object and configure it
	server := &http.Server{
		Addr:    addr,
		Handler: wrappedMux,
	}

	if mode == "production" || https == "true" {
		// HTTPS mode
		certFile := os.Getenv("CERT_FILE")
		keyFile := os.Getenv("KEY_FILE")

		// Use default paths if not specified
		if certFile == "" {
			certFile = "./localhost+2.pem"
		}
		if keyFile == "" {
			keyFile = "./localhost+2-key.pem"
		}

		log.Printf("Starting server in HTTPS mode on port %s\n", port)
		err := server.ListenAndServeTLS(certFile, keyFile)
		if err != nil {
			log.Fatal("Server failed to start:", err)
		}
	} else {
		// HTTP mode (for simple local development)
		log.Printf("Starting server in HTTP mode on port %s\n", port)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Server failed to start:", err)
		}
	}
}
