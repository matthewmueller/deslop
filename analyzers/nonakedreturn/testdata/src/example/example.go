package example

// Bad: naked return with named results.
func divide(a, b int) (result int, err error) {
	if b == 0 {
		return // want "avoid naked return; explicitly return values for readability"
	}
	result = a / b
	return // want "avoid naked return; explicitly return values for readability"
}

// OK: explicit return values.
func add(a, b int) (result int, err error) {
	return a + b, nil
}

// OK: no named returns, bare return not possible with values.
func noop() {
	return
}

// OK: function literal inside should not be checked against outer function.
func outer() (x int) {
	fn := func() {
		return
	}
	_ = fn
	return 0
}
