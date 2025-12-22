package utils

import (
	"context"
	_ "embed"
	"log"
	"strings"

	"github.com/GhostVox/ghostvox.io-backend/internal/database"
	g "github.com/Ghostvox/trie_hard/go"
)

//go:embed restricted_words.txt
var embeddedRestrictedWords string

func parseRestrictedWords(content string) []string {
	lines := strings.Split(content, "\n")
	var words []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		words = append(words, strings.ToLower(line))
	}

	return words
}

// SeedRestrictedWords seeds the database with restricted words from the embedded file
func SeedRestrictedWords(ctx context.Context, db *database.Queries) error {
	words := parseRestrictedWords(embeddedRestrictedWords)

	if len(words) == 0 {
		log.Println("No words to seed")
		return nil
	}

	log.Printf("Seeding %d restricted words to database...", len(words))

	// Try batch insert
	err := db.AddRestrictedWordsBatch(ctx, words)
	if err != nil {
		log.Printf("Batch insert failed, trying one-by-one: %v", err)

		// Fallback to one-by-one
		successCount := 0
		for _, word := range words {
			_, err := db.AddRestrictedWord(ctx, word)
			if err != nil {
				// Silently skip duplicates or errors
				continue
			}
			successCount++
		}
		log.Printf("Successfully seeded %d/%d restricted words", successCount, len(words))
		return nil
	}

	log.Printf("Successfully batch seeded restricted words")
	return nil
}

func NewFilter(db *database.Queries) *g.Trie[string] {
	trie := g.NewTrie[string]()

	// GetAllRestrictedWords now returns []string directly
	restrictedWords, err := db.GetAllRestrictedWords(context.Background())

	var wordsSlice []string

	if err != nil {
		log.Printf("Error fetching from database: %v, using embedded words", err)
		wordsSlice = parseRestrictedWords(embeddedRestrictedWords)

		// Try to seed the database
		seedErr := SeedRestrictedWords(context.Background(), db)
		if seedErr != nil {
			log.Printf("Failed to seed database: %v", seedErr)
		}
	} else if len(restrictedWords) == 0 {
		log.Println("No words in database, seeding from embedded file...")
		wordsSlice = parseRestrictedWords(embeddedRestrictedWords)

		// Seed the database
		seedErr := SeedRestrictedWords(context.Background(), db)
		if seedErr != nil {
			log.Printf("Failed to seed database: %v", seedErr)
		}
	} else {
		// restrictedWords is already []string, just normalize
		wordsSlice = make([]string, 0, len(restrictedWords))
		for _, word := range restrictedWords {
			normalized := strings.ToLower(strings.TrimSpace(word))
			if normalized != "" {
				wordsSlice = append(wordsSlice, normalized)
			}
		}
	}

	if len(wordsSlice) == 0 {
		log.Println("Warning: No restricted words loaded")
		return trie
	}

	trie.AddWordList(&wordsSlice, func(word string) string {
		return strings.Repeat("*", len(word))
	})

	log.Printf("Loaded %d restricted words into filter", len(wordsSlice))
	return trie
}

