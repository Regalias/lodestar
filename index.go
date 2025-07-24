package lodestar

import (
	"fmt"
	"slices"
)

// IndexItems indexes a batch of items, adding them to the immutable radix tree by their tokenized indexes.
// It returns a new Index with the updated index.
// It returns an error if any item fails to index.
func (idx *Index[T]) IndexItems(items []T) (*Index[T], error) {
	if len(items) == 0 {
		return nil, nil
	}

	invertedIndex := make(invertedIndex[T], 0)
	// Build the inverted index
	for _, item := range items {
		if err := idx.addToInvertedIndex(invertedIndex, item); err != nil {
			return nil, fmt.Errorf("failed to index item %v: %w", item, err)
		}
	}

	// Make a new index with the updated immutable radix tree
	newIndex := &Index[T]{
		tokenizer: idx.tokenizer,
	}

	// Index each item in the immutable radix tree
	tx := idx.index.Txn()
	for token, tokenItems := range invertedIndex {
		// Sort the inverted index by descending rank first
		if len(tokenItems) > 1 {
			slices.SortFunc(tokenItems, func(a, b T) int {
				return b.GetRank() - a.GetRank()
			})
			invertedIndex[token] = tokenItems
		}
		// Add the token to the immutable radix tree
		tx.Insert([]byte(token), tokenItems)
	}
	newIndex.index = tx.Commit()

	return newIndex, nil
}

// addToInvertedIndex indexes a single item by tokenizing it and adding the tokens to the inverted index.
func (idx *Index[T]) addToInvertedIndex(invertedIndex invertedIndex[T], item T) error {
	tokens := idx.tokenizer.Tokenize(item)
	if len(tokens) == 0 {
		return fmt.Errorf("no tokens generated for item: %v", item)
	}

	for _, token := range tokens {
		if _, found := invertedIndex[token]; !found {
			invertedIndex[token] = make([]T, 0)
		}
		invertedIndex[token] = append(invertedIndex[token], item)
		// TODO: ensure item is unique?
	}
	return nil
}
