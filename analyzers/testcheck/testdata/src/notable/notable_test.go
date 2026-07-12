package notable

import "testing"

func TestTableInline(t *testing.T) {
	for _, tt := range []struct { // want "avoid table-driven tests; instead either inline the assertions"
		name string
		in   int
	}{
		{"one", 1},
	} {
		_ = tt
	}
}

func TestTableVariable(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"one"},
	}
	for _, tt := range tests { // want "avoid table-driven tests; instead either inline the assertions"
		_ = tt
	}
}

func TestNotATable(t *testing.T) {
	items := []string{"a", "b"}
	for _, item := range items {
		_ = item
	}
}
