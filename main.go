package main

import (
	"github.com/mattmueller/deslop/analyzers/clicheck"
	"github.com/mattmueller/deslop/analyzers/envcheck"
	"github.com/mattmueller/deslop/analyzers/ischeck"
	"github.com/mattmueller/deslop/analyzers/nofmt"
	"github.com/mattmueller/deslop/analyzers/noglobals"
	"github.com/mattmueller/deslop/analyzers/nolog"
	"github.com/mattmueller/deslop/analyzers/notabletest"
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
