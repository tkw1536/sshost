package closer_test

import (
	"errors"
	"fmt"

	"github.com/tkw1536/sshost/internal/pkg/closer"
)

func ExampleStack() {
	stack := closer.NewStack()

	stack.Push(closer.NewCloser(func() error { fmt.Println("second closer"); return nil }))
	stack.Push(closer.NewCloser(func() error { fmt.Println("first closer"); return errors.New("first closer errored") }))

	if err := stack.Close(); err != nil {
		fmt.Printf("error: %q", err)
	}
	// Output:
	// first closer
	// second closer
	// error: "first closer errored"
}
