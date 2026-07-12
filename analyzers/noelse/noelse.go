package noelse

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "noelse",
		Doc:  "reports else blocks; prefer early returns",
		Run:  run,
	}
}

func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		name := pass.Fset.File(f.Pos()).Name()
		if strings.Contains(name, "/go-build/") {
			continue
		}
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		ast.Inspect(f, func(n ast.Node) bool {
			ifStmt, ok := n.(*ast.IfStmt)
			if !ok {
				return true
			}
			if ifStmt.Else != nil {
				pass.Reportf(ifStmt.Else.Pos(), "avoid else blocks; use early returns instead")
			}
			return true
		})
	}
	return nil, nil
}
