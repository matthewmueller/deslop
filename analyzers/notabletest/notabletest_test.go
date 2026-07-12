package notabletest_test

import (
	"testing"

	"github.com/matthewmueller/deslop/analyzers/notabletest"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, notabletest.New(), "example")
}
