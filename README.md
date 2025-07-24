# Lodestar

A generic text search library written in Go

- Immutable indexes
- Support prefix searches
- Sort by rank
- Useful for autocomplete/fast text search on static datasets



## Features

- **Generic Implementation**: Type-safe indexing of any struct that implements the `IndexableItem` interface
- **Prefix Search**: Efficient prefix-based text searching with result ranking
- **Multiple Values**: Index items with multiple searchable values or aliases
- **Result Filtering**: Custom filtering of search results

## Notes

The index performs best when there are generally distinct prefixes in the indexes. The library is designed for in-memory search on static or infrequently changing datasets.

Internally, it uses an immutable radix tree (provided by [hashicorp/go-immutable-radix](http://github.com/hashicorp/go-immutable-radix)) to store and lookup indexes for efficient prefix searching.

Values are indexed by first passing through a tokenizer, which returns a set of unique strings to build the inverted index.

The default tokenizer will generate the following case insensitive indexes, using `the quick_brown-fox (jumps)` as the example value:
- Prefix variations of all word boundaries, split by whitespace and underscores
    - `the quick brown fox jumps`
    - `quick brown fox jumps`
    - `brown fox jumps`
    - `fox jumps`
    - `jumps`
- Variations with and without parenthesis and hyphens.
    - `brown-fox jumps`
    - `brown fox jumps`
    - `(jumps)`
    - `jumps`
    - etc...

This means that a search of `just works` will match all of the following:
- `it just works`
- `it_just_works`
- `it just (works)`
- `it-just-works`

A custom tokenization implementation that implements the `Tokenizer` interface can be provided to support more complex use cases.

## Installation

```bash
go get github.com/regalias/lodestar
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/regalias/lodestar"
)

// Define a struct that implements the IndexableItem interface
type ExampleItem struct {
    Text    string
    Rank    int
    Aliases []string
}

func (e *ExampleItem) GetValuesForIndexing() []string {
    return append([]string{e.Text}, e.Aliases...)
}

func (e *ExampleItem) GetRank() int {
    return e.Rank
}

func (e *ExampleItem) GetID() any {
    return e
}

func main() {
    // Create a new search index with generic type
    index := lodestar.New[*ExampleItem]()
    
    // Create items to index
    items := []*ExampleItem{
        {Text: "apple", Rank: 10, Aliases: []string{"fruit", "red"}},
        {Text: "application", Rank: 15, Aliases: []string{"app", "software"}},
        {Text: "apply", Rank: 8, Aliases: []string{"use", "request"}},
    }
    
    // Index items in batch
    index, err := index.IndexItems(items)
    if err != nil {
        fmt.Println("Error indexing items:", err)
        return
    }
    
    // Perform a prefix search
    results, timing := index.PrefixSearch("app", 10, nil)
    
    // Display results
    fmt.Println("Search results for 'app':")
    for _, result := range results {
        fmt.Printf("- %s (rank: %d)\n", result.Text, result.GetRank())
    }
    
    // Display search timing
    fmt.Printf("Search completed in %v\n", timing.TotalTime)
}
```

## Configuration Options

```go
// Create index with custom tokenizer
index := lodestar.New[*ExampleItem](
    lodestar.WithTokenizer(customTokenizer),
)
```

## API Reference

### Core Types

- `Index[T IndexableItem]`: Generic immutable search index
- `IndexableItem`: Interface for items that can be indexed
- `Tokenizer`: Interface for tokenizing items before indexing
- `ResultFilterFn[T IndexableItem]`: Function type for filtering search results
- `QueryTimingInfo`: Timing information for search operations

### IndexableItem Interface

Any type that implements these methods can be indexed:

```go
type IndexableItem interface {
    // Returns one or more strings to be indexed
    GetValuesForIndexing() []string
    
    // Returns the rank/weight of the item
    GetRank() int
    
    // Returns a unique identifier for deduplication
    GetID() any
}
```

### Main Functions

- `New[T IndexableItem](opts ...Option) *Index[T]`: Create new generic index
- `IndexItems(items []T) (*Index[T], error)`: Index a batch of items
- `PrefixSearch(prefix string, limit int, filterFn ResultFilterFn[T]) ([]T, QueryTimingInfo)`: Search with prefix
- `Get(value string) ([]T, bool)`: Get items by exact match
- `Len() int`: Get number of entries in the index

## Examples

See the [examples](./examples/) directory for more usage patterns.

## Dependencies

- [hashicorp/go-immutable-radix/v2](http://github.com/hashicorp/go-immutable-radix) for the underlying radix tree

## License

See LICENSE file for details.
