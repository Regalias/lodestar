package lodestar

import (
	"testing"
)

func Test_defaultTokenizer_Tokenize(t *testing.T) {

	tr := &DefaultTokenizer{}

	type args struct {
		item IndexableItem
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Single words with aliases",
			args: args{item: &ExampleItem{Text: "apple", Rank: 10, Aliases: []string{"fruit", "red"}}},
			want: []string{"apple", "fruit", "red"},
		},
		{
			name: "Multiple words",
			args: args{item: &ExampleItem{Text: "it just works", Rank: 5, Aliases: []string{}}},
			want: []string{"it just works", "just works", "works"},
		},
		{
			name: "With underscores",
			args: args{item: &ExampleItem{Text: "it_just_works", Rank: 5, Aliases: []string{}}},
			want: []string{"it just works", "just works", "works"},
		},
		{
			name: "With hyphens",
			args: args{item: &ExampleItem{Text: "it-just-works", Rank: 5, Aliases: []string{}}},
			want: []string{"it just works", "just works", "works", "it-just-works"},
		},
		{
			name: "with parentheses",
			args: args{item: &ExampleItem{Text: "it (just) works", Rank: 5, Aliases: []string{}}},
			want: []string{"it just works", "just works", "works", "it (just) works", "(just) works"},
		},
		{
			name: "with everything",
			args: args{item: &ExampleItem{Text: "it (just) works-really!", Rank: 5, Aliases: []string{}}},
			want: []string{
				"it just works really!",
				"just works really!",
				"works really!",
				"really!",

				"it just works-really!",
				"just works-really!",
				"works-really!",
				"it (just) works-really!",
				"(just) works-really!",

				"it (just) works really!",
				"(just) works really!",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tr.Tokenize(tt.args.item)

			// Order insensitive comparison
			gotMap := make(map[string]bool)
			for _, g := range got {
				gotMap[g] = true
			}
			wantMap := make(map[string]bool)
			for _, w := range tt.want {
				wantMap[w] = true
			}

			// Check if all elements in got are in want
			for value := range gotMap {
				if !wantMap[value] {
					t.Errorf("defaultTokenizer.Tokenize() got %v\n\tmissing element %v", got, value)
				}
			}

			// Check if all elements in want are in got
			for value := range wantMap {
				if !gotMap[value] {
					t.Errorf("defaultTokenizer.Tokenize() got %v\n\tmissing element %v", got, value)
				}
			}
		})
	}
}
