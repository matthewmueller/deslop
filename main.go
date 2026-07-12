package main

import (
	"github.com/matthewmueller/deslop/analyzers/clicheck"
	"github.com/matthewmueller/deslop/analyzers/complexity"
	"github.com/matthewmueller/deslop/analyzers/contextfirst"
	"github.com/matthewmueller/deslop/analyzers/envcheck"
	"github.com/matthewmueller/deslop/analyzers/errpath"
	"github.com/matthewmueller/deslop/analyzers/errwrap"
	"github.com/matthewmueller/deslop/analyzers/lowercaseerror"
	"github.com/matthewmueller/deslop/analyzers/muxcheck"
	"github.com/matthewmueller/deslop/analyzers/nocommentedcode"
	"github.com/matthewmueller/deslop/analyzers/noelse"
	"github.com/matthewmueller/deslop/analyzers/nofmt"
	"github.com/matthewmueller/deslop/analyzers/nogetter"
	"github.com/matthewmueller/deslop/analyzers/noglobals"
	"github.com/matthewmueller/deslop/analyzers/noinit"
	"github.com/matthewmueller/deslop/analyzers/nolog"
	"github.com/matthewmueller/deslop/analyzers/nonakedreturn"
	"github.com/matthewmueller/deslop/analyzers/nostutter"
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
		errpath.New(),
		noinit.New(),
		contextfirst.New(),
		nonakedreturn.New(),
		errwrap.New(),
		nocommentedcode.New(),
		noelse.New(),
		nogetter.New(),
		lowercaseerror.New(),
		nostutter.New(),
	}
	analyzers = append(analyzers, testcheck.New()...)
	analyzers = append(analyzers, complexity.New()...)
	multichecker.Main(analyzers...)
}
