package main

import (
	"log"
	"os"

	"github.com/Brent-the-carpenter/ghostvox.io-backend/internal/database"
	"github.com/joho/godotenv"
)
type apiConfig struct {
	db				 *database.Queries
	platform		 string
	port 		   string

}
func main() {
	const (
		port = ":8080"
	)
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	db, err :=
}
