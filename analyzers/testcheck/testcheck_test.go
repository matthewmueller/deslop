package testcheck_test

import (
	"testing"

	"github.com/matthewmueller/deslop/analyzers/testcheck"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestNoTable(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, testcheck.NewNoTable(), "notable")
}

func TestIsCheck(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, testcheck.NewIsCheck(), "ischeck")
}

func TestNaming(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, testcheck.NewNaming(), "naming")
}
