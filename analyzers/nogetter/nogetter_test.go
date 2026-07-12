package nogetter_test

import (
	"testing"

	"github.com/matthewmueller/deslop/analyzers/nogetter"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, nogetter.New(), "example")
}
