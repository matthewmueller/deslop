package lowercaseerror_test

import (
	"testing"

	"github.com/matthewmueller/deslop/analyzers/lowercaseerror"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, lowercaseerror.New(), "example")
}
