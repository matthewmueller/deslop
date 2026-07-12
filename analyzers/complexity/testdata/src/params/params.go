package params

func good(a, b, c, d int) {}

func bad(a, b, c, d, e int) {} // want "function bad has 5 parameters"
