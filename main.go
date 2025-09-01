package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"database/sql"

	"github.com/GhostVox/ghostvox.io-backend/internal/config"
	"github.com/GhostVox/ghostvox.io-backend/internal/cron"
	"github.com/GhostVox/ghostvox.io-backend/internal/database"
	"github.com/GhostVox/ghostvox.io-backend/internal/handlers"
	mw "github.com/GhostVox/ghostvox.io-backend/internal/middleware"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/lib/pq"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

func main() {

	const (
		port = ":8080"
	)

	envConfig, err := LoadEnv()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	//Connect to database
	db, err := sql.Open("postgres", envConfig.DBURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
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
		Platform:          envConfig.Platform,
		Port:              port,
		AccessTokenExp:    envConfig.AccessTokenExp,
		RefreshTokenExp:   envConfig.RefreshTokenExp,
		GhostvoxSecretKey: envConfig.GhostvoxSecretKey,
		Mode:              envConfig.Mode,
		UseHTTPS:          envConfig.UseHTTPS,
		AccessOrigin:      envConfig.AccessOrigin,
		AwsS3Bucket:       envConfig.AWSBucket,
		AwsRegion:         envConfig.AWSRegion,
	}

	//Configure Cron
	CronCFG := cron.NewCronConfig(envConfig.CronCheckExpiredPolls)

	// OAuth2 configuration
	googleOAuthConfig := &oauth2.Config{
		ClientID:     envConfig.GoogleClientID,
		ClientSecret: envConfig.GoogleClientSecret,
		RedirectURL:  envConfig.GoogleRedirectURI,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	githubOAuthConfig := &oauth2.Config{
		ClientID:     envConfig.GithubClientID,
		ClientSecret: envConfig.GithubClientSecret,
		RedirectURL:  envConfig.GithubRedirectURI,
		Scopes: []string{
			"user:email",
			"read:user",
		},
		Endpoint: github.Endpoint,
	}

	credsProvider := credentials.NewStaticCredentialsProvider(envConfig.AWSAccessKeyID, envConfig.AWSSecretAccessKey, "")

	// 1. Load configuration, automatically finding region and credentials
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion(envConfig.AWSRegion), awsconfig.WithCredentialsProvider(credsProvider))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(awsCfg)

	// Initialize handlers

	rootHandler := handlers.NewRootHandler(cfg)
	pollHandler := handlers.NewPollHandler(cfg)
	commentHandler := handlers.NewCommentHandler(cfg)
	voteHandler := handlers.NewVoteHandler(cfg)
	optionHandler := handlers.NewOptionHandler(cfg)
	authHandler := handlers.NewAuthHandler(cfg)
	googleHandler := handlers.NewGoogleHandler(cfg, googleOAuthConfig)
	githubHandler := handlers.NewGithubHandler(cfg, githubOAuthConfig)
	adminHandler := handlers.NewAdminHandler(cfg)
	awsS3Handler := handlers.NewAWSS3Handler(cfg, s3Client)
	userHandler := handlers.NewUserHandler(cfg, awsS3Handler)

	mux := http.NewServeMux()

	wrappedMux := mw.CorsMiddleware(mux, envConfig.AccessOrigin)
	//  start attaching route handlers to cfg.mux
	// Redirects users to the root of the API and returns route endpoints for the API
	mux.HandleFunc("/api/v1/", mw.LoggingMiddleware(rootHandler.HandleRoot)) // in use

	// Polls route ✅
	mux.HandleFunc("GET /api/v1/polls", mw.LoggingMiddleware(pollHandler.GetAllPolls))

	mux.HandleFunc("GET /api/v1/polls/finished", mw.LoggingMiddleware(pollHandler.GetAllFinishedPolls)) // in use
	mux.HandleFunc("GET /api/v1/polls/active", mw.LoggingMiddleware(pollHandler.GetAllActivePolls))     // in use

	mux.HandleFunc("GET /api/v1/polls/recent", mw.LoggingMiddleware(pollHandler.GetRecentPolls)) // in use

	mux.HandleFunc("GET /api/v1/polls/{pollId}", mw.LoggingMiddleware(pollHandler.GetPollByID)) // in use

	mux.HandleFunc("GET /api/v1/polls/{pollId}/comments", mw.LoggingMiddleware(commentHandler.GetAllPollComments))

	mux.HandleFunc("GET /api/v1/users/{userId}/polls", mw.LoggingMiddleware(pollHandler.GetUsersPolls)) // in use

	mux.HandleFunc("PUT /api/v1/polls/{pollId}", mw.LoggingMiddleware(pollHandler.UpdatePoll))

	mux.HandleFunc("POST /api/v1/polls", mw.LoggingMiddleware(pollHandler.CreatePoll))

	mux.HandleFunc("POST /api/v1/polls/{pollId}/vote", mw.LoggingMiddleware(voteHandler.VoteOnPoll))

	mux.HandleFunc("POST /api/v1/polls/{pollId}/comments", mw.LoggingMiddleware(commentHandler.CreatePollComment))

	mux.HandleFunc("DELETE /api/v1/polls/{pollId}/comments/{commentId}", mw.LoggingMiddleware(commentHandler.DeletePollComment))

	mux.HandleFunc("DELETE /api/v1/polls/{pollId}", mw.LoggingMiddleware(pollHandler.DeletePoll))
	// End of poll routes
	// OAuth routes
	mux.HandleFunc("GET /api/v1/auth/google/login", mw.LoggingMiddleware(googleHandler.GoogleLoginHandler))       // in use
	mux.HandleFunc("GET /api/v1/auth/google/callback", mw.LoggingMiddleware(googleHandler.GoogleCallbackHandler)) // in use
	mux.HandleFunc("GET /api/v1/auth/github/login", mw.LoggingMiddleware(githubHandler.GithubLoginHandler))       // in use
	mux.HandleFunc("GET /api/v1/auth/github/callback", mw.LoggingMiddleware(githubHandler.GithubCallbackHandler)) // in use
	// end
	//Auth routes
	mux.HandleFunc("POST /api/v1/auth/login", mw.LoggingMiddleware(authHandler.Login))       // in use
	mux.HandleFunc("POST /api/v1/auth/register", mw.LoggingMiddleware(authHandler.Register)) // in use
	mux.HandleFunc("POST /api/v1/auth/logout", mw.LoggingMiddleware(authHandler.Logout))     // in use
	mux.HandleFunc("POST /api/v1/auth/refresh", mw.LoggingMiddleware(authHandler.Refresh))   // in use

	// Users Private route
	mux.HandleFunc("GET /api/v1/admin/users/{userId}", mw.AdminRole(cfg, mw.LoggingMiddleware(adminHandler.GetUser)).ServeHTTP)

	mux.HandleFunc("GET /api/v1/admin/users", mw.AdminRole(cfg, mw.LoggingMiddleware(adminHandler.GetAllUsers)).ServeHTTP)

	// User public routes
	mux.HandleFunc("GET /api/v1/users/stats", mw.LoggingMiddleware(userHandler.GetUserStats))
	mux.HandleFunc("PUT /api/v1/users/profile", mw.LoggingMiddleware(userHandler.UpdateUser))
	mux.HandleFunc("PUT /api/v1/users/profile/avatar", mw.LoggingMiddleware(awsS3Handler.UpdateUserAvatar))
	mux.HandleFunc("POST /api/v1/users/username", mw.LoggingMiddleware(userHandler.AddUserName))
	mux.HandleFunc("DELETE /api/v1/users/", mw.LoggingMiddleware(userHandler.DeleteUser))

	// End of users routes

	// Options Routes

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

	if envConfig.Mode == "production" || envConfig.UseHTTPS == "true" {
		// HTTPS mode
		certFile := os.Getenv("TLS_CERT_PATH")
		keyFile := os.Getenv("TLS_KEY_PATH")

		// Use default paths if not specified
		if envConfig.CertFile == "" {
			envConfig.CertFile = "./localhost+2.pem"
		}
		if envConfig.KeyFile == "" {
			envConfig.KeyFile = "./localhost+2-key.pem"
		}
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)

		log.Println("Starting cron jobs")
		go func() {
			cronCtx := context.Background()
			CronCFG.StartCronJobs(cronCtx, cfg)
		}()
		go func() {
			log.Printf("Starting server in HTTPS mode on port %s\n", port)
			err := server.ListenAndServeTLS(certFile, keyFile)
			if err != nil {
				log.Fatal("Server failed to start:", err)
			}
		}()
		<-stop
		fmt.Println("Shuting down server and cron jobs")
		serverCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		err := server.Shutdown(serverCtx)
		if err != nil {
			log.Fatal("Failed to shutdown server:", err)
		}
		CronCFG.StopJobs()
		fmt.Println("Server and cron jobs stopped")
	} else {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)

		log.Println("Starting cron jobs")
		go func() {
			cronCtx := context.Background()
			CronCFG.StartCronJobs(cronCtx, cfg)
		}()

		// HTTP mode (for simple local development)
		go func() {
			log.Printf("Starting server in HTTP mode on port %s\n", port)
			err := server.ListenAndServe()
			if err != nil {
				log.Fatal("Server failed to start:", err)
			}
		}()

		<-stop
		log.Println("Stopping server")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Server shutdown failed:", err)
		}
		CronCFG.StopJobs()
		fmt.Println("Server and cron jobs shutdown.")
	}
}
