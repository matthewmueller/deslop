package example

import "os"

func bad() {
	_ = os.MkdirAll("/tmp/foo", 0755) // want "do not discard errors"
}

func badMultiReturn() {
	_, _ = os.Open("/tmp/foo") // want "do not discard errors"
}

func badInline() {
	if f, _ := os.Open("/tmp/foo"); f != nil { // want "do not discard errors"
		f.Close()
	}
}

// OK: error is captured.
func good() {
	err := os.MkdirAll("/tmp/foo", 0755)
	_ = err
}

// OK: first return discarded, error captured.
func goodMulti() {
	_, err := os.Open("/tmp/foo")
	_ = err
}

// OK: not a function call.
func goodNonCall() {
	m := map[string]int{"a": 1}
	_, _ = m["a"], m["b"]
}

func goodNonError() {
	_ = func() string { return "hello" }() // OK: not an error
}
