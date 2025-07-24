package lodestar

import "time"

// IndexableItem defines the interface for items that can be indexed.
type IndexableItem interface {
	// GetValuesForIndexing returns one or more strings associated with the item.
	// All strings will be tokenized into a unique set of tokens for indexing
	GetValuesForIndexing() []string

	// GetRank returns the rank of the index item.
	GetRank() int

	// Returns any hashable value that uniquely identifies the item, used for deduplication.
	// A basic implementation could return a pointer of the item itself
	GetID() any
}

// Tokenizer defines the interface for tokenizing items before indexing
type Tokenizer interface {
	// Tokenize tokenizes the values of an IndexableItem into a set of unique tokens.
	Tokenize(item IndexableItem) []string

	// TokenizeString tokenizes a single string value into a set of unique tokens.
	// TokenizeString(value string) []string

	// NormalizeString normalizes a string for tokenization (e.g. lowercasing, removing punctuation).
	// This is also applied to all queries before searching
	NormalizeString(value string) string
}

// ResultFilterFn is a function to filter results based on custom logic.
// It takes the normalized original query, matching token, and item as parameters and returns true if the item
// should be included in aggregation.
// Duplicate items are automatically filtered out by their GetID() value.
// If nil, all items are included.
type ResultFilterFn[T IndexableItem] func(normalizedQuery string, token string, item T) bool

// QueryTimingInfo holds timing information for a search query execution.
type QueryTimingInfo struct {
	InitTime        time.Duration
	SeekTime        time.Duration
	AggregationTime time.Duration
	TotalTime       time.Duration
}
