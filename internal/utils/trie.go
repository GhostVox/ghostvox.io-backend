package utils

type TrieNode struct {
	children    map[rune]*TrieNode
	isEndOfWord bool
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		children:    make(map[rune]*TrieNode),
		isEndOfWord: false,
	}
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: NewTrieNode(),
	}
}

func (t *Trie) Insert(word string) {
	currNode := t.root
	for _, char := range word {
		if _, ok := currNode.children[char]; !ok {
			currNode.children[char] = NewTrieNode()
		}
		currNode = currNode.children[char]
	}
	currNode.isEndOfWord = true
}

func (t *Trie) Search(word string) bool {
	currNode := t.root
	for _, char := range word {
		if _, ok := currNode.children[char]; !ok {
			return false
		}
		currNode = currNode.children[char]
	}
	return currNode.isEndOfWord
}

func (t *Trie) StartsWith(prefix string) bool {
	currNode := t.root
	for _, char := range prefix {
		if _, ok := currNode.children[char]; !ok {
			return false
		}
		if currNode.isEndOfWord {
			return false
		}
		currNode = currNode.children[char]
	}
	return true
}
