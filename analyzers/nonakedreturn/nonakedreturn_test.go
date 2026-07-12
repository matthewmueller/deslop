package nonakedreturn_test

import (
	"testing"

	"github.com/matthewmueller/deslop/analyzers/nonakedreturn"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, nonakedreturn.New(), "example")
}
