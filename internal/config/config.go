package config

import (
	"database/sql"
	"time"

	"github.com/GhostVox/ghostvox.io-backend/internal/database"
)

type APIConfig struct {
	DB                *sql.DB
	Queries           *database.Queries
	Platform          string
	Port              string
	AccessTokenExp    time.Duration
	RefreshTokenExp   time.Duration
	GhostvoxSecretKey string
	Mode              string
	UseHTTPS          string
	AccessOrigin      string
	AwsS3Bucket       string
	AwsRegion         string
	DOMAIN            string
}
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
