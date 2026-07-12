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

// Allowed: unexported vars
var db = connect()
var client interface{}
var config = &myError{}
var lookup = map[string]bool{"a": true}

func connect() interface{} { return nil }

// Not allowed: exported vars (not Err*)
var DB = connect()         // want `avoid package-level var "DB"`
var DefaultClient = "http" // want `avoid package-level var "DefaultClient"`
