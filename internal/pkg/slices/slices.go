// Package slices contains utility functions for slices.
package slices

import "errors"

// ErrNoValue indicates that no value in slice was contained in valid.
var ErrNoValue = errors.New("no valid value found")

// Filter filters slice by elements that only occur in valid.
// Does not re-allocate, and invalidates memory used by s.
//
// When slice is nil, never returns an error.
// When slice becomes empty, returns ErrNoValue.
func Filter[T comparable](slice []T, valid []T) ([]T, error) {
	// special case: 0-size slice or valid
	if len(slice) == 0 || len(valid) == 0 {
		return slice, nil
	}

	// cache which elements exist
	cache := make(map[T]struct{}, len(valid))
	for _, v := range valid {
		cache[v] = struct{}{}
	}

	// filter s according to the cache
	result := slice[:0]
	for _, element := range slice {
		if _, ok := cache[element]; ok {
			result = append(result, element)
		}
	}

	if len(result) == 0 {
		return slice, ErrNoValue
	}

	return result, nil
}

// Combine creates a new slice that contains elements from all slices passed in order.
// Elements are copied into a new slice, and the original slices are left unchanged.
func Combine[T any](slices ...[]T) (result []T) {
	for _, s := range slices {
		result = append(result, s...)
	}
	return
}
