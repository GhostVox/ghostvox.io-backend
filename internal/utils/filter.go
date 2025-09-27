package utils

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/GhostVox/ghostvox.io-backend/internal/database"
	g "github.com/Ghostvox/trie_hard/go"
)

func NewFilter(db *database.Queries) *g.Trie[string] {
	trie := g.NewTrie[string]()

	restrictedWords, err := db.GetAllRestrictedWords(context.Background())
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Error fetching restricted words: %v", err)
		}
		// Initialize with empty trie rather than panicking
		return trie
	}

	// Assuming restrictedWords is already a []string or can be converted correctly
	// If restrictedWords needs type conversion, do it appropriately here
	restrictedWordsSlice, ok := restrictedWords.([]string)
	if !ok {
		log.Printf("Error: restrictedWords is not of type []string")
		return trie
	}

	trie.AddWordList(&restrictedWordsSlice, func(word string) string {
		return strings.Repeat("*", len(word))
	})

	return trie
}
