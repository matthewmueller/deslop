package testcheck

import (
	"strings"

	"golang.org/x/tools/go/analysis"
)

func hasTestFiles(pass *analysis.Pass) bool {
	for _, f := range pass.Files {
		name := pass.Fset.File(f.Pos()).Name()
		if strings.HasSuffix(name, "_test.go") {
			return true
		}
	}
	return false
}
