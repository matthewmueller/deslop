package ischeck_test

import (
	"testing"

	"github.com/matthewmueller/deslop/analyzers/ischeck"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ischeck.New(), "example")
}
