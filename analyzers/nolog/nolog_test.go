package nolog_test

import (
	"testing"

	"github.com/matthewmueller/deslop/analyzers/nolog"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, nolog.New(), "example")
}
