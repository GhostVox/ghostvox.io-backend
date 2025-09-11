package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
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
	AWSRegion             string
	AWSBucket             string
	AWSAccessKeyID        string
	AWSSecretAccessKey    string
	IPRateLimit           rate.Limit
	IPRateBurst           int
}

// LoadEnv loads environment variables and returns a config struct
func LoadEnv() (*EnvConfig, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		// Continue anyway, as env vars might be set directly in the environment
	}

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

	https := getRequiredEnv("USE_HTTPS")
	var dbURL string
	if https == "true" {
		dbURL = "postgres://" + DB_USER + ":" + DB_PASSWORD + "@" + DB_HOST + ":" + DB_PORT + "/" + DB_NAME + "?sslmode=disable"
	} else {
		dbURL = getRequiredEnv("DATABASE_URL")
	}
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

	awsRegion := getRequiredEnv("AWS_REGION")
	awsBucket := getRequiredEnv("AWS_S3_BUCKET")
	awsAccessKeyID := getRequiredEnv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := getRequiredEnv("AWS_SECRET_ACCESS_KEY")

	mode := getRequiredEnv("MODE")
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
	certFile := os.Getenv("CERT_FILE")
	if certFile == "" {
		certFile = "./localhost+2.pem"
	}

	keyFile := os.Getenv("KEY_FILE")
	if keyFile == "" {
		keyFile = "./localhost+2-key.pem"
	}

	ipRateLimitStr := os.Getenv("IP_RATE_LIMIT")
	if ipRateLimitStr == "" {
		ipRateLimitStr = "10"
	}
	rateLimit, err := strconv.Atoi(ipRateLimitStr)
	if err != nil {
		log.Fatalf("Invalid IP rate limit: %v", err)
	}

	ipRateLimit := rate.Limit(rateLimit)

	ipRateBurstStr := os.Getenv("IP_RATE_BURST")
	if ipRateBurstStr == "" {
		ipRateBurstStr = "100"
	}
	ipRateBurst, err := strconv.Atoi(ipRateBurstStr)
	if err != nil {
		log.Fatalf("Invalid IP rate burst: %v", err)
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
		AWSRegion:             awsRegion,
		AWSBucket:             awsBucket,
		AWSAccessKeyID:        awsAccessKeyID,
		AWSSecretAccessKey:    awsSecretAccessKey,
		IPRateLimit:           ipRateLimit,
		IPRateBurst:           ipRateBurst,
	}, nil
}
