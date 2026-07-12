package muxcheck_test

import (
	"testing"

	"github.com/matthewmueller/deslop/analyzers/muxcheck"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, muxcheck.New(), "example")
}
