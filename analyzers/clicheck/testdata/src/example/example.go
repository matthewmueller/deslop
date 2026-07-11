package example

import (
	"flag" // want `do not use the "flag" package directly`
	"os"
)

func Bad() {
	_ = os.Args // want "do not access os.Args directly"
	flag.Parse()
}
