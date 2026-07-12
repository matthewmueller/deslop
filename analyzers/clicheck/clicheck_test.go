package clicheck_test

import (
	"testing"

	"github.com/mattmueller/deslop/analyzers/clicheck"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, clicheck.New(), "example")
}
