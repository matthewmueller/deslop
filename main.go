package main

import (
	"github.com/matthewmueller/deslop/analyzers/clicheck"
	"github.com/matthewmueller/deslop/analyzers/envcheck"
	"github.com/matthewmueller/deslop/analyzers/ischeck"
	"github.com/matthewmueller/deslop/analyzers/nofmt"
	"github.com/matthewmueller/deslop/analyzers/noglobals"
	"github.com/matthewmueller/deslop/analyzers/nolog"
	"github.com/matthewmueller/deslop/analyzers/notabletest"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		nofmt.New(),
		envcheck.New(),
		clicheck.New(),
		notabletest.New(),
		ischeck.New(),
		noglobals.New(),
		nolog.New(),
	)
}
