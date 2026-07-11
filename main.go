package main

import (
	"github.com/mattmueller/deslop/analyzers/nofmt"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		nofmt.Analyzer,
	)
}
