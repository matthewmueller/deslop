package noerrdiscard_test

import (
	"testing"

	"github.com/matthewmueller/deslop/analyzers/noerrdiscard"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, noerrdiscard.New(), "example")
}
