package lodestar

import (
	"fmt"
	"testing"
)

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

var testItems = []*ExampleItem{
	{Text: "apple", Rank: 10, Aliases: []string{"fruit", "red"}},
	{Text: "banana", Rank: 12, Aliases: []string{"yellow", "fruit"}},
	{Text: "application", Rank: 15, Aliases: []string{"app", "software"}},
	{Text: "apply", Rank: 8, Aliases: []string{"use", "request"}},
	{Text: "approach", Rank: 5, Aliases: []string{"method", "way"}},
}

func setupEmptyIndex() *Index[*ExampleItem] {
	return New[*ExampleItem]()
}

func setupIndexWithItems(items []*ExampleItem) *Index[*ExampleItem] {
	index := setupEmptyIndex()
	index, _ = index.IndexItems(items)
	return index
}

func TestNew(t *testing.T) {
	index := New[*ExampleItem]()
	if index == nil {
		t.Fatal("New() returned nil")
	}
	if index.Len() != 0 {
		t.Errorf("New index should be empty, got size %d", index.Len())
	}
}

func TestInsert(t *testing.T) {
	index := setupEmptyIndex()

	// Test single insert
	newIndex, err := index.IndexItems([]*ExampleItem{
		{Text: "test", Rank: 10, Aliases: []string{"example"}},
	})
	if err != nil {
		t.Fatalf("Failed to insert item: %v", err)
	}

	if newIndex.Len() != 2 { // 1 value + 1 alias
		t.Errorf("Expected size 2 after insert, got %d", newIndex.Len())
	}

	// Original index should be unchanged (immutable)
	if index.Len() != 0 {
		t.Errorf("Original index should remain unchanged, got size %d", index.Len())
	}
}

func BenchmarkInsert(b *testing.B) {
	index := setupEmptyIndex()

	for i := 0; b.Loop(); i++ {
		index, _ = index.IndexItems([]*ExampleItem{
			{Text: fmt.Sprintf("item%d", i), Rank: i, Aliases: []string{"example"}},
		})
	}
}
