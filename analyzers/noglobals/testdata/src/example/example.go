package example

import "errors"

// Allowed: Err* sentinels
var ErrNotFound = errors.New("not found")

// Allowed: compile-time interface check
var _ error = (*myError)(nil)

type myError struct{}

func (e *myError) Error() string { return "" }

// Not allowed: globals
var db interface{} // want `avoid package-level var "db"`

var defaultClient = "http://localhost" // want `avoid package-level var "defaultClient"`

// Allowed: constants are fine
const defaultPort = 8080
