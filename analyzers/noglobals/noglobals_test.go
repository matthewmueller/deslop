package noglobals_test

import (
	"testing"

	"github.com/matthewmueller/deslop/analyzers/noglobals"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, noglobals.New(), "example")
}
