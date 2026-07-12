package nesting

func good() {
	if true {
		if true {
			if true {
				// depth 3 — OK
			}
		}
	}
}

func bad() { // want "function bad has nesting depth 4"
	if true {
		if true {
			if true {
				if true {
					// depth 4
				}
			}
		}
	}
}
