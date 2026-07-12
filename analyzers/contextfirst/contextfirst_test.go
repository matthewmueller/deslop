package contextfirst_test

import (
	"testing"

	"github.com/matthewmueller/deslop/analyzers/contextfirst"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, contextfirst.New(), "example")
}
