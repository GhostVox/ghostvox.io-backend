package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"

	"github.com/GhostVox/ghostvox.io-backend/internal/database"
	"github.com/GhostVox/ghostvox.io-backend/internal/handlers"
	mw "github.com/GhostVox/ghostvox.io-backend/internal/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db       *database.Queries
	platform string
	port     string
	mux      *http.ServeMux
}

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
	cfg := apiConfig{
		db:       dbConnection,
		platform: platform,
		port:     port,
	}
	rootHandler := handlers.NewRootHandler(cfg.db)
	pollHandler := handlers.NewPollHandler(cfg.db)
	mux := http.NewServeMux()
	//  start attaching route handlers to cfg.mux

	// Redirects users to the root of the API and returns route endpoints for the API
	mux.HandleFunc("/api/v1/", mw.LoggingMiddleware(rootHandler.HandleRoot))

	// Polls route
	mux.HandleFunc("GET /api/v1/polls", mw.LoggingMiddleware(pollHandler.GetAllPolls))

	mux.HandleFunc("GET /api/v1/polls/{pollId}", mw.LoggingMiddleware(pollHandler.GetPoll))

	mux.HandleFunc("PUT /api/v1/polls/{pollId}", mw.LoggingMiddleware(pollHandler.UpdatePoll))

	mux.HandleFunc("POST /api/v1/polls", mw.LoggingMiddleware(pollHandler.CreatePoll))

	mux.HandleFunc("DELETE /api/v1/polls/{pollId}", mw.LoggingMiddleware(pollHandler.DeletePoll))
	// End of poll routes

	mux.HandleFunc("POST /api/v1/polls/{pollId}/options", func(w http.ResponseWriter, r *http.Request) {})

	mux.HandleFunc("GET /api/v1/polls/{pollId}/options", func(w http.ResponseWriter, r *http.Request) {})

	mux.HandleFunc("PUT /api/v1/polls/{pollId}/options/{optionId}", func(w http.ResponseWriter, r *http.Request) {})

	mux.HandleFunc("/api/v1/polls/{pollId}/votes/{voteId}", func(w http.ResponseWriter, r *http.Request) {})

	mux.HandleFunc("/api/v1/polls/{pollId}/votes/{voteId}/results", func(w http.ResponseWriter, r *http.Request) {})

	mux.HandleFunc("GET/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "http://localhost:8080/api/v1", http.StatusFound)
			return
		}

		http.NotFound(w, r)
	}))

	// Create a pointer to the http.Server object and configure it
	server := &http.Server{
		Addr:    cfg.port,
		Handler: mux,
	}

	// Start the server
	log.Printf("Server running on http://localhost%s", cfg.port)
	log.Fatal(server.ListenAndServe())

}
