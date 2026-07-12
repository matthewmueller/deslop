package errpath_test

import (
	"testing"

	"github.com/matthewmueller/deslop/analyzers/errpath"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, errpath.New(), "example")
}
