package lodestar

// Option represents a configuration option for the Index.
type Option func(*Config)

// Config holds configuration for the Index.
type Config struct {
	Tokenizer Tokenizer
}

// WithTokenizer returns an Option that sets the tokenizer for the index.
func WithTokenizer(tokenizer Tokenizer) Option {
	return func(c *Config) {
		c.Tokenizer = tokenizer
	}
}
