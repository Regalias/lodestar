package lodestar

import (
	"container/heap"
	"log/slog"
	"math"
	"time"

	"github.com/regalias/lodestar/internal/utils"
)

// result represents a search result with associated weight.
type result[T IndexableItem] struct {
	Value T
	Rank  int
}

// resultHeap implements heap.Interface for sorting search results by rank.
type resultHeap[T IndexableItem] []result[T]

func (h resultHeap[T]) Len() int           { return len(h) }
func (h resultHeap[T]) Less(i, j int) bool { return h[i].Rank < h[j].Rank }
func (h resultHeap[T]) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *resultHeap[T]) Push(x any) {
	*h = append(*h, x.(result[T]))
}
func (h *resultHeap[T]) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

// Search performs a prefix search and returns results sorted by rank (descending) up to the specified limit.
// The prefix is normalized before searching.
// If a filter function is provided, it will be applied to each item before including it in the results.
// The results are deduplicated based on the item's GetID() value.
func (idx *Index[T]) PrefixSearch(prefix string, limit int, filterFn ResultFilterFn[T]) ([]T, QueryTimingInfo) {

	t0 := time.Now()
	prefix = idx.tokenizer.NormalizeString(prefix)

	if prefix == "" {
		return nil, QueryTimingInfo{}
	}

	// Default to no limit
	if limit <= 0 {
		limit = math.MaxInt
	}

	// Set for deduplication
	seen := make(map[any]struct{}, limit)

	// Min-heap to get top K
	minHeap := resultHeap[T]{}
	heap.Init(&minHeap)

	iter := idx.index.Root().Iterator()

	t1 := time.Now()

	iter.SeekPrefix([]byte(prefix))

	t2 := time.Now()

	for token, items, ok := iter.Next(); ok; token, items, ok = iter.Next() {
		tokenStr := string(token)
		if len(items) == 0 {
			slog.Warn("Got empty items for token", "token", tokenStr)
			continue
		}

		for _, item := range items {

			// Skip duplicates
			if _, exists := seen[item.GetID()]; exists {
				continue
			}
			seen[item.GetID()] = struct{}{}

			// Apply the filter function if provided
			if filterFn != nil && !filterFn(prefix, tokenStr, item) {
				continue
			}

			if len(minHeap) < limit {
				// Not enough items, just add it
				heap.Push(&minHeap, result[T]{Value: item, Rank: item.GetRank()})
			} else if item.GetRank() > minHeap[0].Rank {
				// The current item has a higher rank than the lowest in the heap, replace it
				heap.Pop(&minHeap)
				heap.Push(&minHeap, result[T]{Value: item, Rank: item.GetRank()})
			} else {
				// Current item is not better than the worst item in the heap, stop here
				// The rest of the items will have lower rank since they are already sorted by descending rank
				break
			}
		}
	}

	t3 := time.Now()

	// Get the results from min-heap in ascending order of rank
	results := make([]T, 0, len(minHeap))
	for minHeap.Len() > 0 {
		result := heap.Pop(&minHeap).(result[T])
		results = append(results, result.Value)
	}

	// Reverse the results to get them in descending order of rank
	utils.ReverseSliceInPlace(results)
	t4 := time.Now()

	return results, QueryTimingInfo{
		InitTime:        t1.Sub(t0),
		SeekTime:        t2.Sub(t1),
		AggregationTime: t3.Sub(t2),
		TotalTime:       t4.Sub(t0),
	}
}

// Get retrieves an item by exact match
// The value is normalized before searching the indexes
func (idx *Index[T]) Get(value string) (item []T, found bool) {
	value = idx.tokenizer.NormalizeString(value)
	if value == "" {
		return item, false
	}
	root := idx.index.Root()
	return root.Get([]byte(value))
}
