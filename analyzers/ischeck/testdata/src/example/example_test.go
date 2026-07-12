package example

import "testing"

func TestBadFatal(t *testing.T) {
	t.Fatal("something failed") // want "avoid t.Fatal for assertions"
}

func TestBadErrorf(t *testing.T) {
	t.Errorf("got %d, want %d", 1, 2) // want "avoid t.Errorf for assertions"
}

// Helper functions are allowed to use t.Fatal.
func setup(t *testing.T) {
	t.Helper()
	t.Fatal("setup failed")
}

func TestUsesHelper(t *testing.T) {
	setup(t)
}
