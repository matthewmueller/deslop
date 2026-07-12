package noelse_test

import (
	"testing"

	"github.com/matthewmueller/deslop/analyzers/noelse"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, noelse.New(), "example")
}
