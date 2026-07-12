package nostutter_test

import (
	"testing"

	"github.com/matthewmueller/deslop/analyzers/nostutter"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, nostutter.New(), "user")
}
