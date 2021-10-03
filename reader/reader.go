// Package reader implements reading of configuration settings
package reader

// Reader reads configuration settings for aliases
type Reader interface {
	// Get returns a single value for the given alias and key
	// Non-existing keys may return either the empty string or a sensible default.
	Get(alias, key string) (string, error)

	// GetAll returns all values for a given alias and key.
	// Non-existing keys should return a zero-length slice.
	GetAll(alias, key string) ([]string, error)
}
