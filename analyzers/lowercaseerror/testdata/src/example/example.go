package example

import "fmt"

func bad() error {
	err := fmt.Errorf("other")
	return fmt.Errorf("Parsing failed: %w", err) // want `error strings should start with lowercase; got "Parsing failed: %w"`
}

func good() error {
	err := fmt.Errorf("other")
	return fmt.Errorf("parsing failed: %w", err)
}

func dynamic(msg string) error {
	return fmt.Errorf(msg)
}
