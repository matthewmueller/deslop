package envcheck_test

import (
	"testing"

	"github.com/mattmueller/deslop/analyzers/envcheck"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, envcheck.Analyzer, "example")
}
