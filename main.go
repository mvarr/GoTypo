package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// ----- Trie Implementation -----

// TrieNode represents a node in our Trie.
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

// NewTrieNode creates a new Trie node.
func NewTrieNode() *TrieNode {
	return &TrieNode{children: make(map[rune]*TrieNode)}
}

// Trie holds our root TrieNode.
type Trie struct {
	root *TrieNode
}

// NewTrie initializes a new Trie.
func NewTrie() *Trie {
	return &Trie{root: NewTrieNode()}
}

// Insert adds a word into the Trie.
func (t *Trie) Insert(word string) {
	node := t.root
	for _, ch := range word {
		if _, exists := node.children[ch]; !exists {
			node.children[ch] = NewTrieNode()
		}
		node = node.children[ch]
	}
	node.isEnd = true
}

// Search checks if a word exists in the Trie.
func (t *Trie) Search(word string) bool {
	node := t.root
	for _, ch := range word {
		if next, ok := node.children[ch]; ok {
			node = next
		} else {
			return false
		}
	}
	return node.isEnd
}

// ----- Loading Dictionary from File -----

// LoadWords reads words from a file (one word per line), inserts them into the Trie,
// and returns a slice of words for suggestion purposes.
func LoadWords(filename string, trie *Trie) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			words = append(words, word)
			trie.Insert(word)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return words, nil
}

// ----- Levenshtein Distance Algorithm -----

// Levenshtein calculates the edit distance between two strings.
func Levenshtein(a, b string) int {
	m := len(a)
	n := len(b)
	if m == 0 {
		return n
	}
	if n == 0 {
		return m
	}

	// Initialize a 2D slice.
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	// Base cases.
	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}

	// Fill dp table.
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			dp[i][j] = min(
				dp[i-1][j]+1,      // deletion
				dp[i][j-1]+1,      // insertion
				dp[i-1][j-1]+cost, // substitution
			)
		}
	}

	return dp[m][n]
}

func min(a, b, c int) int {
	if a < b && a < c {
		return a
	} else if b < c {
		return b
	}
	return c
}

// ----- Suggestion Algorithm -----

// SuggestWords returns words from the dictionary with an edit distance <= maxDistance.
// It sorts the suggestions by edit distance and limits the returned list to maxSuggestions.
func SuggestWords(word string, dictionary []string, maxDistance int, maxSuggestions int) []string {
	var suggestions []string
	for _, dictWord := range dictionary {
		d := Levenshtein(word, dictWord)
		if d <= maxDistance {
			suggestions = append(suggestions, dictWord)
		}
	}
	// Sort suggestions by edit distance (lower distance first)
	sort.Slice(suggestions, func(i, j int) bool {
		return Levenshtein(word, suggestions[i]) < Levenshtein(word, suggestions[j])
	})
	// Limit the number of suggestions returned
	if len(suggestions) > maxSuggestions {
		suggestions = suggestions[:maxSuggestions]
	}
	return suggestions
}

// BestSuggestion picks the suggestion with the smallest edit distance.
func BestSuggestion(word string, suggestions []string) string {
	if len(suggestions) == 0 {
		return ""
	}
	best := suggestions[0]
	bestDistance := Levenshtein(word, best)
	for _, s := range suggestions {
		d := Levenshtein(word, s)
		if d < bestDistance {
			best = s
			bestDistance = d
		}
	}
	return best
}

// ----- Main CLI Application -----

func main() {
	trie := NewTrie()

	// Load dictionary from words.txt.
	dictionary, err := LoadWords("words.txt", trie)
	if err != nil {
		fmt.Println("Error loading dictionary:", err)
		return
	}

	// Interactive mode.
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("GoTypo - Enter a word to check (type 'exit' to quit):")
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		input = strings.TrimSpace(input)
		input = strings.ToLower(input)
		if input == "exit" {
			break
		}
		if trie.Search(input) {
			fmt.Printf("'%s' is correct.\n", input)
		} else {
			fmt.Printf("'%s' is not in the dictionary.\n", input)
			// Use a maxDistance of 2 and limit suggestions to 10 words.
			suggestions := SuggestWords(input, dictionary, 2, 10)
			if len(suggestions) == 0 {
				fmt.Println("No suggestions found.")
			} else {
				fmt.Println("Suggestions:")
				for _, s := range suggestions {
					fmt.Println(" -", s)
				}
				best := BestSuggestion(input, suggestions)
				if best != "" {
					fmt.Printf("Best suggestion: %s\n", best)
				}
			}
		}
	}
}
