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
	mode := os.Getenv("MODE")
	if mode != "" {
		log.Fatalf("Error loading MODE: %v", err)
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
		DB:                dbConnection,
		Platform:          platform,
		Port:              port,
		AccessTokenExp:    accesstokenExpDur,
		RefreshTokenExp:   refreshTokenExpDur,
		GhostvoxSecretKey: ghostvoxSecretKey,
		Mode:              mode,
	}
	// Initialize handlers
	rootHandler := handlers.NewRootHandler(cfg)
	pollHandler := handlers.NewPollHandler(cfg.DB)
	voteHandler := handlers.NewVoteHandler(cfg.DB)
	optionHandler := handlers.NewOptionHandler(cfg)
	userHandler := handlers.NewUserHandler(cfg)
	authHandler := handlers.NewAuthHandler(cfg)
	googleHandler := handlers.NewGoogleHandler(cfg)

	mux := http.NewServeMux()
	wrappedMux := mw.CorsMiddleware(mux)
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
	mux.Handle("/api/v1/auth/google/login", mw.LoggingMiddleware(googleHandler.GoogleLoginHandler))
	mux.Handle("/api/v1/auth/google/callback", mw.LoggingMiddleware(googleHandler.GoogleCallbackHandler))
	// end
	mux.HandleFunc("POST /api/v1/auth/google", mw.LoggingMiddleware(authHandler.GoogleOAuth))
	//Auth routes
	mux.HandleFunc("POST /api/v1/auth/login", mw.LoggingMiddleware(authHandler.Login))
	mux.HandleFunc("POST /api/v1/auth/register", mw.LoggingMiddleware(authHandler.Register))
	mux.HandleFunc("POST /api/v1/auth/logout", mw.LoggingMiddleware(authHandler.Logout))
	mux.HandleFunc("POST /api/v1/auth/refresh", mw.LoggingMiddleware(authHandler.Refresh))

	// Users Private route
	mux.HandleFunc("GET /api/v1/admin/users/{userId}", mw.AdminRole(cfg, mw.LoggingMiddleware(userHandler.GetUser)).ServeHTTP)

	mux.HandleFunc("GET /api/v1/admin/users", mw.AdminRole(cfg, mw.LoggingMiddleware(userHandler.GetAllUsers)).ServeHTTP)

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

	// Create a pointer to the http.Server object and configure it
	server := &http.Server{
		Addr:    cfg.Port,
		Handler: wrappedMux,
	}

	// Start the server
	log.Printf("Server running on http://localhost%s", cfg.Port)
	log.Fatal(server.ListenAndServe())

}
