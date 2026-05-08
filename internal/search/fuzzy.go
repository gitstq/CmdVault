// Package search provides fuzzy search capabilities for CmdVault
package search

import (
	"sort"

	"github.com/sahilm/fuzzy"
	"github.com/cmdvault/cmdvault/internal/storage"
)

// SearchResult represents a search result with match information
type SearchResult struct {
	Command    storage.Command
	Score      int
	MatchedIdx []int
}

// FuzzySearcher provides fuzzy search functionality
type FuzzySearcher struct {
	commands []storage.Command
}

// NewFuzzySearcher creates a new fuzzy searcher
func NewFuzzySearcher(commands []storage.Command) *FuzzySearcher {
	return &FuzzySearcher{commands: commands}
}

// Search performs fuzzy search on commands
func (fs *FuzzySearcher) Search(query string, limit int) []SearchResult {
	if query == "" {
		// Return top commands by execution count
		results := make([]SearchResult, 0, min(limit, len(fs.commands)))
		for i := 0; i < limit && i < len(fs.commands); i++ {
			results = append(results, SearchResult{
				Command: fs.commands[i],
				Score:   fs.commands[i].ExecCount * 100,
			})
		}
		return results
	}

	// Create string slice for fuzzy matching
	names := make([]string, len(fs.commands))
	for i, cmd := range fs.commands {
		names[i] = cmd.Command
	}

	matches := fuzzy.Find(query, names)

	// Sort by score
	sort.Sort(sort.Reverse(matches))

	// Convert to results
	results := make([]SearchResult, 0, min(limit, len(matches)))
	for i := 0; i < limit && i < len(matches); i++ {
		match := matches[i]
		results = append(results, SearchResult{
			Command:    fs.commands[match.Index],
			Score:      match.Score,
			MatchedIdx: match.MatchedIndexes,
		})
	}

	return results
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
