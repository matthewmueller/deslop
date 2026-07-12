package noglobals_test

import (
	"testing"

	"github.com/mattmueller/deslop/analyzers/noglobals"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, noglobals.Analyzer, "example")
}
