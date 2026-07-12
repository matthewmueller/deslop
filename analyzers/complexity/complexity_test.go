package complexity_test

import (
	"testing"

	"github.com/matthewmueller/deslop/analyzers/complexity"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestNesting(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, complexity.NewNesting(), "nesting")
}

func TestParams(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, complexity.NewParams(), "params")
}

func TestNaming(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, complexity.NewNaming(), "naming")
}

func TestLineLength(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, complexity.NewLineLength(), "linelength")
}

func TestFuncLength(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, complexity.NewFuncLength(), "funclength")
}
