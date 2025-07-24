package lodestar

import (
	iradix "github.com/hashicorp/go-immutable-radix/v2"
)

type invertedIndex[T IndexableItem] map[string][]T

// Index represents an immutable text search index.
type Index[T IndexableItem] struct {
	index     *iradix.Tree[[]T]
	tokenizer Tokenizer
}

// New creates a new empty Index. If no Tokenizer is provided, it uses the default tokenizer,
// which does case insensitive indexing of all word boundaries, split by whitespace, underscores, and hyphens.
func New[T IndexableItem](opts ...Option) *Index[T] {
	config := &Config{}
	for _, opt := range opts {
		opt(config)
	}

	// Use a default tokenizer if none provided
	if config.Tokenizer == nil {
		config.Tokenizer = &DefaultTokenizer{}
	}

	return &Index[T]{
		index:     iradix.New[[]T](),
		tokenizer: config.Tokenizer,
	}
}

// Len returns the number of items in the index radix tree
func (idx *Index[T]) Len() int {
	return idx.index.Len()
}
