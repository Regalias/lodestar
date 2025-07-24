package lodestar

import (
	"fmt"
	"math/rand"
	"testing"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

var testItems2 = []*ExampleItem{
	{Text: "the quick brown fox", Rank: 20, Aliases: []string{"the (quick) brown fox"}},
	{Text: "the slow brown fox", Rank: 18, Aliases: []string{}},
	{Text: "the very-slow brown fox", Rank: 22, Aliases: []string{}},
	{Text: "the quick brown-dog", Rank: 25, Aliases: []string{}},
}

// TODO: make this more robust

func TestSearch(t *testing.T) {
	index := setupIndexWithItems(testItems)
	results, _ := index.PrefixSearch("app", 0, nil)
	if len(results) != 4 {
		t.Errorf("Expected 4 results, got %d", len(results))
	}

	index = setupIndexWithItems(testItems2)
	results, _ = index.PrefixSearch("brown fox", 0, nil)
	if len(results) != 3 {
		t.Errorf("Expected 3 results for 'brown fox', got %d", len(results))
	}

	results, _ = index.PrefixSearch("dog", 0, nil)
	if len(results) != 1 {
		t.Errorf("Expected 1 result for 'dog', got %d", len(results))
	}

	results, _ = index.PrefixSearch("slow", 0, nil)
	if len(results) != 2 {
		t.Errorf("Expected 2 results for 'slow', got %d", len(results))
	}
}

func BenchmarkSearchBadCase(b *testing.B) {
	index := setupEmptyIndex()

	// Setup a large number of items with similar prefixes
	items := make([]*ExampleItem, 0, 250000)
	for i := range 250000 {
		items = append(items, &ExampleItem{
			Text:    fmt.Sprintf("item%d", i),
			Rank:    i,
			Aliases: []string{"example"},
		})
	}
	index, _ = index.IndexItems(items)
	b.Logf("Index size: %d", index.Len())

	for b.Loop() {
		index.PrefixSearch("item", 100, nil)
	}
}

func BenchmarkSearchGoodCase(b *testing.B) {

	// This case is technically the worst for memory usage

	index := setupEmptyIndex()

	toQuery := []string{}

	items := make([]*ExampleItem, 0, 250000)

	// Setup a large number of items with distinct prefixes
	for i := range 250000 {

		id1 := generateRandomString(5)
		id2 := generateRandomString(5)
		id3 := generateRandomString(5)

		if i == 0 {
			toQuery = append(toQuery, id1, id2, id3)
		}

		id := fmt.Sprintf("%s_%s_%s", id1, id2, id3)
		items = append(items, &ExampleItem{
			Text:    id,
			Rank:    i,
			Aliases: []string{},
		})
	}

	index, _ = index.IndexItems(items)
	b.Logf("Index size: %d", index.Len())

	for i := 0; b.Loop(); i++ {
		index.PrefixSearch(toQuery[i%3], 100, nil)
	}
}

func BenchmarkSearchMixedCase(b *testing.B) {
	index := setupEmptyIndex()

	id1 := ""
	id2 := ""
	id3 := ""

	toQuery := []string{}

	items := make([]*ExampleItem, 0, 250000)

	// Setup a large number of items with mixed prefixes
	for i := range 250000 {

		if i%250 == 0 {
			id1 = generateRandomString(4)
		}
		if i%333 == 0 {
			id2 = generateRandomString(4)
		}
		if i%713 == 0 {
			id3 = generateRandomString(4)
		}

		if i == 0 {
			toQuery = append(toQuery, id1, id2, id3)
		}

		id := fmt.Sprintf("%s_%s_%s", id1, id2, id3)
		items = append(items, &ExampleItem{
			Text:    id,
			Rank:    i,
			Aliases: []string{},
		})
	}

	index, _ = index.IndexItems(items)
	b.Logf("Index size: %d", index.Len())

	for i := 0; b.Loop(); i++ {
		index.PrefixSearch(toQuery[i%3], 100, nil)
	}
}
