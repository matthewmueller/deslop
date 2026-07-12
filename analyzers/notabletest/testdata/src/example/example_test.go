package example

import "testing"

func TestTableInline(t *testing.T) {
	for _, tt := range []struct { // want "avoid table-driven tests"
		name string
		in   int
		out  int
	}{
		{"one", 1, 1},
		{"two", 2, 4},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.in*tt.in != tt.out {
				t.Fatal("mismatch")
			}
		})
	}
}

func TestTableVariable(t *testing.T) {
	tests := []struct {
		name string
		in   int
	}{
		{"one", 1},
		{"two", 2},
	}
	for _, tt := range tests { // want "avoid table-driven tests"
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.in
		})
	}
}

func TestNotATable(t *testing.T) {
	items := []string{"a", "b", "c"}
	for _, item := range items {
		_ = item
	}
}
