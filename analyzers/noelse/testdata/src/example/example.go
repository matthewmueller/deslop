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
	} else { // want "avoid else blocks; use early returns instead"
		return "non-positive"
	}
}

func badElseIf(x int) string {
	if x > 0 {
		return "positive"
	} else if x == 0 { // want "avoid else blocks; use early returns instead"
		return "zero"
	}
	return "negative"
}
