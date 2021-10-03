package sshost

import (
	"errors"
	"fmt"
)

func ExampleClosableStack() {
	closer := NewClosableStack()

	closer.Push(NewCloser(func() error { fmt.Println("second closer"); return nil }))
	closer.Push(NewCloser(func() error { fmt.Println("first closer"); return errors.New("first closer errored") }))

	if err := closer.Close(); err != nil {
		fmt.Printf("error: %q", err)
	}
	// Output:
	// first closer
	// second closer
	// error: "first closer errored"
}
