package errwrap_test

import (
	"testing"

	"github.com/matthewmueller/deslop/analyzers/errwrap"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, errwrap.New(), "example")
}
