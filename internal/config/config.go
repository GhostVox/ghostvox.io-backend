package config

import (
	"time"

	"github.com/GhostVox/ghostvox.io-backend/internal/database"
)

type APIConfig struct {
	DB                *database.Queries
	Platform          string
	Port              string
	AccessTokenExp    time.Duration
	RefreshTokenExp   time.Duration
	GhostvoxSecretKey string
}
