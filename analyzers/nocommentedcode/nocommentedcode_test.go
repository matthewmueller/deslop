package nocommentedcode_test

import (
	"testing"

	"github.com/matthewmueller/deslop/analyzers/nocommentedcode"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, nocommentedcode.New(), "example")
}
