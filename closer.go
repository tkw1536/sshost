package sshost

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

// CloseableStack is a stack-like data structure that contains multiple internal Closers.
// It is safe to be used concurrently by multiple goroutines.
type ClosableStack struct {
	m       sync.RWMutex
	closers []Closer
}

// NewClosable stack creates a new CloseableStack
func NewClosableStack(closers ...Closer) *ClosableStack {
	stack := &ClosableStack{}
	stack.Push(closers...)
	return stack
}

// Push adds several closers to the internal stack, in order.
// stack may not be nil.
//
// If a current call to Close() is in progress, blocks until such a call is finished.
func (stack *ClosableStack) Push(closers ...Closer) {
	if stack == nil {
		panic("CloseableStack.Push: stack is nil")
	}

	stack.m.Lock()
	defer stack.m.Unlock()

	stack.closers = append(stack.closers, closers...)
}

// PushStack is like Push, except that it takes a CloseableStack as argument
func (stack *ClosableStack) PushStack(other *ClosableStack) {
	stack.Push(stack.closers...)
}

// Reset resets this stack to an empty state.
// Waits until all calls to Close() have finished.
//
// When stack is nil, does nothing.
func (stack *ClosableStack) Reset() {
	if stack == nil {
		return
	}

	stack.m.Lock()
	defer stack.m.Unlock()

	stack.closers = nil
}

// Close calls all closers in this stack in LIFO order.
// It is safe to make multiple calls to Close() simultaniously, either will run independently.
//
// Returns the first non-nil error that was returned by any closer.
//
// The implementation ensures that all closers are called, even if one panics.
// If such a panic occurs the panic will not be recovered.
func (stack *ClosableStack) Close() (err error) {
	if stack == nil {
		return nil
	}

	stack.m.RLock()
	defer stack.m.RUnlock()

	for _, closer := range stack.closers {
		defer func(closer Closer) {
			if e := closer.Close(); e != nil && err == nil {
				err = e
			}
		}(closer)
	}

	return
}
