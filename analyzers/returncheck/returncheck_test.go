package returncheck_test

import (
	"testing"

	"github.com/mattmueller/deslop/analyzers/returncheck"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, returncheck.Analyzer, "example")
}
