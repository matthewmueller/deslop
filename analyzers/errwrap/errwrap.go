package errwrap

import (
	"go/ast"
	"go/types"
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
	errIface := types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

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
			if !hasErrorArg(pass, call.Args[1:], errIface) {
				return true
			}
			pass.Reportf(call.Pos(), "use %%w in fmt.Errorf to wrap errors instead of %%v or %%s")
			return true
		})
	}
	return nil, nil
}

func hasErrorArg(pass *analysis.Pass, args []ast.Expr, errIface *types.Interface) bool {
	return slices.ContainsFunc(args, func(expr ast.Expr) bool {
		t := pass.TypesInfo.TypeOf(expr)
		if t == nil {
			return false
		}
		return types.Implements(t, errIface)
	})
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
