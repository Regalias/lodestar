// Package lodestar provides a generic immutable text search library supporting
// prefix search and weights for autocomplete functionality.
//
// This library is built on top of hashicorp's immutable radix tree
// implementation and provides a type-safe API for building searchable
// text indexes with weighted results.
//
// Any type that implements the IndexableItem interface can be indexed:
//
//	type IndexableItem interface {
//		// Returns one or more strings to be indexed
//		GetValuesForIndexing() []string
//
//		// Returns the rank/weight of the item
//		GetRank() int
//
//		// Returns a unique identifier for deduplication
//		GetID() any
//	}
//
// Basic usage:
//
//	// Define a struct that implements IndexableItem
//	type ExampleItem struct {
//		Text    string
//		Rank    int
//		Aliases []string
//	}
//
//	func (e *ExampleItem) GetValuesForIndexing() []string {
//		return append([]string{e.Text}, e.Aliases...)
//	}
//
//	func (e *ExampleItem) GetRank() int {
//		return e.Rank
//	}
//
//	func (e *ExampleItem) GetID() any {
//		return e
//	}
//
//	// Create a new index with generic type
//	index := lodestar.New[*ExampleItem]()
//
//	// Index items in batch
//	items := []*ExampleItem{
//		{Text: "apple", Rank: 10, Aliases: []string{"fruit"}},
//		{Text: "application", Rank: 15, Aliases: []string{"app"}},
//	}
//	index, err := index.IndexItems(items)
//
//	// Perform a prefix search
//	results, timing := index.PrefixSearch("app", 10, nil)
//	for _, result := range results {
//		fmt.Printf("%s (rank: %d)\n", result.Text, result.GetRank())
//	}
package lodestar
