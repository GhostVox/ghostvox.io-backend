package main

import (
	"log"
	"os"
	"time"
	//"github.com/joho/godotenv"
)

// EnvConfig holds all environment configuration
type EnvConfig struct {
	DBURL                 string
	Platform              string
	GhostvoxSecretKey     string
	AccessTokenExp        time.Duration
	RefreshTokenExp       time.Duration
	GoogleClientID        string
	GoogleClientSecret    string
	GoogleRedirectURI     string
	GithubClientID        string
	GithubClientSecret    string
	GithubRedirectURI     string
	Mode                  string
	UseHTTPS              string
	AccessOrigin          string
	CronCheckExpiredPolls string
	CertFile              string
	KeyFile               string
}

// LoadEnv loads environment variables and returns a config struct
func LoadEnv() (*EnvConfig, error) {
	// Load environment variables from .env file
	//	err := godotenv.Load()
	//	if err != nil {
	//		log.Println("Error loading .env file")
	//		// Continue anyway, as env vars might be set directly in the environment
	//	}
	//
	// Define a helper function for required env vars
	getRequiredEnv := func(key string) string {
		val := os.Getenv(key)
		if val == "" {
			log.Fatalf("%s must be set", key)
		}
		return val
	}
	// Get the  DB connection URL parts

	DB_HOST := getRequiredEnv("DB_HOST")
	DB_PORT := getRequiredEnv("DB_PORT")
	DB_NAME := getRequiredEnv("DB_NAME")
	DB_USER := getRequiredEnv("DB_USER")
	DB_PASSWORD := getRequiredEnv("DB_PASSWORD")

	dbURL := "postgres://" + DB_USER + ":" + DB_PASSWORD + "@" + DB_HOST + ":" + DB_PORT + "/" + DB_NAME + "?sslmode=disable"
	// Get all required env vars
	platform := getRequiredEnv("PLATFORM")
	secretKey := getRequiredEnv("GHOSTVOX_SECRET_KEY")
	accessTokenExpStr := getRequiredEnv("ACCESS_TOKEN_EXPIRES")
	refreshTokenExpStr := getRequiredEnv("REFRESH_TOKEN_EXPIRES")
	googleClientID := getRequiredEnv("GOOGLE_CLIENT_ID")
	googleClientSecret := getRequiredEnv("GOOGLE_CLIENT_SECRET")
	googleRedirectURI := getRequiredEnv("GOOGLE_REDIRECT_URI")
	githubClientID := getRequiredEnv("GITHUB_CLIENT_ID")
	githubClientSecret := getRequiredEnv("GITHUB_CLIENT_SECRET")
	githubRedirectURI := getRequiredEnv("GITHUB_REDIRECT_URI")
	mode := getRequiredEnv("MODE")
	https := getRequiredEnv("USE_HTTPS")
	accessOrigin := getRequiredEnv("ACCESS_ORIGIN")
	cronCheckExpiredPolls := getRequiredEnv("CRON_CHECK_FOR_EXPIRED_POLLS")

	// Parse durations
	accessTokenExp, err := time.ParseDuration(accessTokenExpStr)
	if err != nil {
		log.Fatalf("Invalid access token expiration: %v", err)
	}

	refreshTokenExp, err := time.ParseDuration(refreshTokenExpStr)
	if err != nil {
		log.Fatalf("Invalid refresh token expiration: %v", err)
	}

	// Get optional env vars with defaults
	certFile := os.Getenv("TLS_CERT_PATH")
	if certFile == "" {
		certFile = "./localhost+2.pem"
	}

	keyFile := os.Getenv("TLS_KEY_PATH")
	if keyFile == "" {
		keyFile = "./localhost+2-key.pem"
	}

	return &EnvConfig{
		DBURL:                 dbURL,
		Platform:              platform,
		GhostvoxSecretKey:     secretKey,
		AccessTokenExp:        accessTokenExp,
		RefreshTokenExp:       refreshTokenExp,
		GoogleClientID:        googleClientID,
		GoogleClientSecret:    googleClientSecret,
		GoogleRedirectURI:     googleRedirectURI,
		GithubClientID:        githubClientID,
		GithubClientSecret:    githubClientSecret,
		GithubRedirectURI:     githubRedirectURI,
		Mode:                  mode,
		UseHTTPS:              https,
		AccessOrigin:          accessOrigin,
		CronCheckExpiredPolls: cronCheckExpiredPolls,
		CertFile:              certFile,
		KeyFile:               keyFile,
	}, nil
}
