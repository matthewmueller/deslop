package aliased

import ctx "context"

// Bad: context.Context is not first, even with aliased import.
func Bad(name string, c ctx.Context) {} // want "context.Context should be the first parameter of Bad"

// OK: context.Context is first with aliased import.
func Good(c ctx.Context, name string) {}
