package main

import (
	"fmt"

	"github.com/regalias/lodestar"
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

func main() {
	// Create a new index
	index := lodestar.New[*ExampleItem]()

	items := []*ExampleItem{
		{Text: "apple", Rank: 10, Aliases: []string{"fruit", "red"}},
		{Text: "banana", Rank: 12, Aliases: []string{"yellow", "fruit"}},
		{Text: "application", Rank: 15, Aliases: []string{"app", "software"}},
		{Text: "apply", Rank: 8, Aliases: []string{"use", "request"}},
		{Text: "approach", Rank: 5, Aliases: []string{"method", "way"}},
	}

	// Add some entries, returning a new immutable Index
	index, err := index.IndexItems(items)
	if err != nil {
		fmt.Println("Error indexing items:", err)
		return
	}

	// Search for entries starting with "app"
	results, _ := index.PrefixSearch("app", 0, nil)

	fmt.Println("Search results for 'app':")
	for _, result := range results {
		fmt.Printf("- %v (weight: %d)\n", result, result.GetRank())
	}

	results, found := index.Get("apply")
	if !found {
		fmt.Println("No results found for 'apply'")
		return
	}
	fmt.Println("Get results for 'apply':")
	for _, result := range results {
		fmt.Printf("- %v (weight: %d)\n", result, result.GetRank())
	}
}
