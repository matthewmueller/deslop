package nofmt_test

import (
	"testing"

	"github.com/matthewmueller/deslop/analyzers/nofmt"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, nofmt.New(), "example")
}
