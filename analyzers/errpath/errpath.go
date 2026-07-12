package errpath

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// New returns an analyzer that flags "if err == nil" patterns.
func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "errpath",
		Doc:      "reports error handling patterns that put the happy path on the wrong side",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      run,
	}
}

func run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.IfStmt)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		ifStmt := n.(*ast.IfStmt)
		if isErrEqualsNil(ifStmt.Cond) {
			pass.Reportf(ifStmt.Pos(), "avoid \"if err == nil\"; invert to \"if err != nil { return ... }\" and keep the happy path on the left. If the nil path is genuinely shorter, use \"nil == err\" to suppress this warning.")
		}
	})

	return nil, nil
}

// isErrEqualsNil checks for "err == nil" (but NOT "nil == err").
func isErrEqualsNil(expr ast.Expr) bool {
	binExpr, ok := expr.(*ast.BinaryExpr)
	if !ok || binExpr.Op != token.EQL {
		return false
	}
	if isErrIdent(binExpr.X) && isNilIdent(binExpr.Y) {
		return true
	}
	// "nil == err" is allowed as an escape hatch.
	return false
}

func isErrIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "err"
}

func isNilIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "nil"
}
