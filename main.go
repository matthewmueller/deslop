package main

import (
	"github.com/mattmueller/deslop/analyzers/clicheck"
	"github.com/mattmueller/deslop/analyzers/envcheck"
	"github.com/mattmueller/deslop/analyzers/nofmt"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		nofmt.Analyzer,
		envcheck.Analyzer,
		clicheck.Analyzer,
	)
}
