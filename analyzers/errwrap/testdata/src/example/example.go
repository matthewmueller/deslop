package example

import (
	"errors"
	"fmt"
)

func good(err error) error {
	return fmt.Errorf("failed to connect: %w", err)
}

func bad(err error) error {
	return fmt.Errorf("failed to connect: %v", err) // want "use %w in fmt.Errorf to wrap errors instead of %v or %s"
}

// OK: creating a new error with context data, no error arg.
func newError(name string) error {
	return fmt.Errorf("invalid name: %s", name)
}

// OK: no args beyond format string.
func simpleError() error {
	return fmt.Errorf("something went wrong")
}

// OK: non-literal format string, can't check statically.
func dynamic(msg string, err error) error {
	return fmt.Errorf(msg, err)
}

var _ = errors.New("not checked")
