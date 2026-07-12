package errwrap

import (
	"go/ast"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "errwrap",
		Doc:  "reports fmt.Errorf calls that pass an error without wrapping with %w",
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
			if len(call.Args) < 2 {
				return true
			}
			lit, ok := call.Args[0].(*ast.BasicLit)
			if !ok {
				return true
			}
			if strings.Contains(lit.Value, "%w") {
				return true
			}
			if !hasErrorArg(call.Args[1:]) {
				return true
			}
			pass.Reportf(call.Pos(), "use %%w in fmt.Errorf to wrap errors instead of %%v or %%s")
			return true
		})
	}
	return nil, nil
}

func hasErrorArg(args []ast.Expr) bool {
	return slices.ContainsFunc(args, isErrorExpr)
}

func isErrorExpr(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		name := e.Name
		return name == "err" || strings.HasPrefix(name, "err") || strings.HasSuffix(name, "Err") || strings.HasSuffix(name, "err")
	case *ast.CallExpr:
		sel, ok := e.Fun.(*ast.SelectorExpr)
		if !ok {
			return false
		}
		return sel.Sel.Name == "Error"
	case *ast.SelectorExpr:
		name := e.Sel.Name
		return name == "Err" || strings.HasPrefix(name, "err") || strings.HasSuffix(name, "Err")
	}
	return false
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
