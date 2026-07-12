package example

import "context"

// OK: context.Context is first parameter.
func Good(ctx context.Context, name string) {}

// OK: no context.Context at all.
func NoContext(name string, age int) {}

// Bad: context.Context is not first.
func Bad(name string, ctx context.Context) {} // want "context.Context should be the first parameter of Bad"

// OK: unexported function, not checked.
func bad(name string, ctx context.Context) {}

type Server struct{}

// OK: method with context first.
func (s *Server) Handle(ctx context.Context, req string) {}

// Bad: method with context not first.
func (s *Server) Process(req string, ctx context.Context) {} // want "context.Context should be the first parameter of Process"
