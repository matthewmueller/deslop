package testcheck

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// NewNaming reports test functions with underscores in their names.
func NewNaming() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "testnaming",
		Doc:      "reports test functions with underscores (e.g., Test_Something); use TestSomething instead",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      runNaming,
	}
}

func runNaming(pass *analysis.Pass) (any, error) {
	if !hasTestFiles(pass) {
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		fn := n.(*ast.FuncDecl)
		pos := pass.Fset.Position(fn.Pos())
		if !strings.HasSuffix(pos.Filename, "_test.go") {
			return
		}
		if !strings.HasPrefix(fn.Name.Name, "Test") {
			return
		}
		// Flag Test_Something but not TestSomething.
		name := fn.Name.Name
		after := strings.TrimPrefix(name, "Test")
		if strings.HasPrefix(after, "_") {
			pass.Reportf(fn.Name.Pos(), "avoid underscores in test names; use %q instead of %q", "Test"+strings.TrimPrefix(after, "_"), name)
		}
	})

	return nil, nil
}
