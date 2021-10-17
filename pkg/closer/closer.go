// Package closer provides utilities for handling io.Closers
package closer

import (
	"io"
	"sync"
)

// Closer is an alias for io.Closer.
type Closer = io.Closer

type closerFunc func() error

func (c closerFunc) Close() error {
	return c()
}

// NewCloser creates a new Closer from a function.
func NewCloser(close func() error) Closer {
	return closerFunc(close)
}

// Stack is a stack-like data structure that contains multiple internal Closers.
// It is safe to be used concurrently by multiple goroutines.
type Stack struct {
	m       sync.RWMutex
	closers []Closer
}
