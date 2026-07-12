package errwrap

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "errwrap",
		Doc:  "reports fmt.Errorf calls without %w verb",
		Run:  run,
	}
}

func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			if !isFmtErrorf(call) {
				return true
			}
			if len(call.Args) == 0 {
				return true
			}
			lit, ok := call.Args[0].(*ast.BasicLit)
			if !ok {
				// Can't statically check non-literal format strings.
				return true
			}
			// lit.Value includes quotes, e.g. `"failed: %w"`
			if !strings.Contains(lit.Value, "%w") {
				pass.Reportf(call.Pos(), "use %%w in fmt.Errorf to wrap errors; got fmt.Errorf(%s, ...) without %%w", lit.Value)
			}
			return true
		})
	}
	return nil, nil
}

func isFmtErrorf(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == "fmt" && sel.Sel.Name == "Errorf"
}
