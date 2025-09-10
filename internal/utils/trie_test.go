package utils_test

import (
	"testing"

	"github.com/GhostVox/ghostvox.io-backend/internal/utils"
)

func TestTrie(t *testing.T) {
	trie := utils.NewTrie()
	trie.Insert("apple")
	trie.Insert("app")
	trie.Insert("banana")

	if !trie.Search("apple") {
		t.Errorf("Expected 'apple' to be found")
	}

	if !trie.Search("app") {
		t.Errorf("Expected 'app' to be found")
	}

	if !trie.Search("banana") {
		t.Errorf("Expected 'banana' to be found")
	}

	if trie.Search("orange") {
		t.Errorf("Expected 'orange' not to be found")
	}

	if !trie.StartsWith("app") {
		t.Errorf("Expected 'app' to start with 'app'")
	}

	if !trie.StartsWith("ban") {
		t.Errorf("Expected 'banana' to start with 'ban'")
	}

	if trie.StartsWith("orange") {
		t.Errorf("Expected 'orange' not to start with 'orange'")
	}
}
