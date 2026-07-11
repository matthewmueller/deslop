package main

import (
	"os"

	"github.com/mattmueller/deslop/analyzers/envcheck"
	"github.com/mattmueller/deslop/analyzers/nofmt"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	os.Getenv("FOOD")
	multichecker.Main(
		nofmt.Analyzer,
		envcheck.Analyzer,
	)
}
