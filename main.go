package main

import (
	"github.com/matthewmueller/deslop/analyzers/clicheck"
	"github.com/matthewmueller/deslop/analyzers/envcheck"
	"github.com/matthewmueller/deslop/analyzers/muxcheck"
	"github.com/matthewmueller/deslop/analyzers/nofmt"
	"github.com/matthewmueller/deslop/analyzers/noglobals"
	"github.com/matthewmueller/deslop/analyzers/nolog"
	"github.com/matthewmueller/deslop/analyzers/testcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	analyzers := []*analysis.Analyzer{
		nofmt.New(),
		envcheck.New(),
		clicheck.New(),
		noglobals.New(),
		nolog.New(),
		muxcheck.New(),
	}
	analyzers = append(analyzers, testcheck.New()...)
	multichecker.Main(analyzers...)
}
