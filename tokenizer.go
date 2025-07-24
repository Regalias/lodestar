package lodestar

import (
	"log/slog"
	"strings"
)

// DefaultTokenizer implements the Tokenizer interface with basic tokenization logic.
// The string is normalized by lowercasing and trimming whitespace.
// Underscores are replaced with spaces for terms longer than 3 characters.
// It tokenizes strings by splitting on whitespace and computing prefix combinations with and without hyphens.
type DefaultTokenizer struct{}

func (t *DefaultTokenizer) Tokenize(item IndexableItem) []string {
	var allTokens []string
	for _, value := range item.GetValuesForIndexing() {
		tokens := t.tokenizeString(value)
		if tokens != nil {
			allTokens = mergeUniqueTokens(allTokens, tokens)
		} else {
			slog.Warn("Tokenization for item value returned an empty set", "value", value, "item", item)
		}
	}
	return allTokens
}

func (t *DefaultTokenizer) tokenizeString(value string) []string {
	// Normalize the term
	normalizedValue := t.NormalizeString(value)
	terms := splitOnWhitespace(normalizedValue)
	tokens := computePrefixCombinations(terms)
	if len(tokens) == 0 {
		return nil
	}

	// If the value has dashes, also compute this variation
	if strings.Contains(value, "-") && len(value) > 3 {
		dashlessString := strings.ReplaceAll(normalizedValue, "-", " ")
		dashlessTokens := computePrefixCombinations(splitOnWhitespace(dashlessString))
		tokens = mergeUniqueTokens(tokens, dashlessTokens)
	}

	return computeTokenVariations(tokens)
}

// NormalizeString normalizes a term by converting it to lowercase and removing leading/trailing whitespace.
// If the term is longer than 3 characters, it also replaces underscores with spaces.
func (t *DefaultTokenizer) NormalizeString(value string) string {
	normalized := strings.ToLower(value)
	normalized = strings.TrimSpace(normalized)
	if len(normalized) > 3 {
		normalized = strings.ReplaceAll(normalized, "_", " ")
	}
	return normalized
}

// computeTokenVariations computes variations of tokens with special characters that may appear in
// words removed, like bracket.
// This allows us to support searching for values with and without brackets, e.g.:
// "(hello world)" when searching for "hello world" and "(hello world)"
// The order of tokens is not preserved. Any added variations will be unique.
func computeTokenVariations(tokens []string) []string {
	// Use a set to store unique tokens and their variations.
	variationSet := make(map[string]struct{}, len(tokens)*2) // Pre-allocate for variations
	// Replacer for removing all bracket types efficiently.
	bracketReplacer := strings.NewReplacer("(", "", ")", "", "[", "", "]", "", "{", "", "}", "")

	for _, token := range tokens {
		// Always add the original token.
		variationSet[token] = struct{}{}

		// Ignore short tokens for further variations.
		if len(token) < 4 {
			continue
		}
		// Generate variations only if special characters are present.
		if strings.ContainsAny(token, "()[]{}") {
			// Variation with brackets removed.
			noBrackets := bracketReplacer.Replace(token)
			if noBrackets != token {
				variationSet[noBrackets] = struct{}{}
			}
		}
	}

	if len(variationSet) == len(tokens) {
		return tokens
	}

	// Convert the set back to a slice.
	result := make([]string, 0, len(variationSet))
	for v := range variationSet {
		result = append(result, v)
	}
	return result
}

// computePrefixCombinations computes prefix combinations of tokens.
// This is useful for prefix searches where we want to match any prefix of a token.
func computePrefixCombinations(tokens []string) []string {
	if len(tokens) == 0 {
		return nil
	}
	prefixes := make([]string, 0, len(tokens))
	for curr_token_pos := range tokens {
		prefixes = append(prefixes, strings.Join(tokens[curr_token_pos:], " "))
	}
	return prefixes
}

func splitOnWhitespace(normalizedString string) []string {
	// Split by whitespace and filter out empty tokens
	tokens := make([]string, 0)
	for _, token := range strings.Fields(normalizedString) {
		if token != "" {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func mergeUniqueTokens(tokens ...[]string) []string {
	// Use a set to ensure uniqueness
	tokenSet := make(map[string]struct{})
	for _, tokenList := range tokens {
		for _, token := range tokenList {
			tokenSet[token] = struct{}{}
		}
	}

	// Convert the map keys back to a slice
	result := make([]string, 0, len(tokenSet))
	for token := range tokenSet {
		result = append(result, token)
	}

	// TODO: sort for deterministic output?
	return result
}
