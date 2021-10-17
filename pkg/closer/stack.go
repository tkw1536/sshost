package closer

// NewStack creates a new stack with the provided closers
func NewStack(closers ...Closer) *Stack {
	stack := &Stack{}
	stack.Push(closers...)
	return stack
}

// Push adds several closers to the internal stack, in order.
// stack may not be nil.
//
// If a current call to Close() is in progress, blocks until such a call is finished.
func (stack *Stack) Push(closers ...Closer) {
	if stack == nil {
		panic("CloseableStack.Push: stack is nil")
	}

	stack.m.Lock()
	defer stack.m.Unlock()

	stack.closers = append(stack.closers, closers...)
}

// PushStack is like Push, except that it takes a Stack as argument
func (stack *Stack) PushStack(other *Stack) {
	stack.Push(stack.closers...)
}

// Reset resets this stack to an empty state.
// Waits until all calls to Close() have finished.
//
// When stack is nil, does nothing.
func (stack *Stack) Reset() {
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
func (stack *Stack) Close() (err error) {
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
