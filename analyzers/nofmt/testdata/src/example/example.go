package example

import "fmt"

func Hello() {
	fmt.Println("hello") // want "avoid fmt.Println; use structured logging instead"
}
