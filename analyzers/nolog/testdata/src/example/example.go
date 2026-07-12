package example

import "log" // want `do not use the stdlib "log" package`

func Bad() {
	log.Println("hello") // want "do not use log.Println"
	log.Fatalf("err: %v", "oops") // want "do not use log.Fatalf"
}
