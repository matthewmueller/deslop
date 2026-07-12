package example

func good(x int) string {
	if x > 0 {
		return "positive"
	}
	return "non-positive"
}

func bad(x int) string {
	if x > 0 {
		return "positive"
	} else { // want "avoid else blocks; use an early return or extract a separate function instead"
		return "non-positive"
	}
}

// Allowed: else if that returns.
func goodElseIf(x int) string {
	if x > 0 {
		return "positive"
	} else if x == 0 {
		return "zero"
	}
	return "negative"
}

// Bad: else if that doesn't return.
func badElseIf(x int) string {
	if x > 0 {
		return "positive"
	} else if x == 0 { // want "avoid else blocks; use an early return or extract a separate function instead"
		_ = "zero"
	}
	return "negative"
}
