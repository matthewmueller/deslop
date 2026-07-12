package example

import (
	"errors"
	"fmt"
)

func good(err error) error {
	return fmt.Errorf("failed to connect: %w", err)
}

func bad(err error) error {
	return fmt.Errorf("failed to connect: %v", err) // want `use %w in fmt.Errorf to wrap errors; got fmt.Errorf\("failed to connect: %v", \.\.\.\) without %w`
}

func badNoVerb() error {
	return fmt.Errorf("something went wrong") // want `use %w in fmt.Errorf to wrap errors; got fmt.Errorf\("something went wrong", \.\.\.\) without %w`
}

func dynamic(msg string, err error) error {
	// OK: non-literal format string, can't check statically.
	return fmt.Errorf(msg, err)
}

var _ = errors.New("not checked")
