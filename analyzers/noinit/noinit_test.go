package noinit_test

import (
	"testing"

	"github.com/matthewmueller/deslop/analyzers/noinit"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, noinit.New(), "example")
}
