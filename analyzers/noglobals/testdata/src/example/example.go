package example

import "errors"

// Allowed: Err* sentinels
var ErrNotFound = errors.New("not found")

// Allowed: compile-time interface check
var _ error = (*myError)(nil)

type myError struct{}

func (e *myError) Error() string { return "" }

// Allowed: constants
const defaultPort = 8080

// Not allowed: function call result assigned to a var
var db = connect() // want `avoid package-level var "db"`

func connect() interface{} { return nil }

// Not allowed: bare uninitialized var
var client interface{} // want `avoid package-level var "client"`

// Not allowed: composite literals are still globals
var config = &myError{}            // want `avoid package-level var "config"`
var lookup = map[string]bool{"a": true} // want `avoid package-level var "lookup"`
