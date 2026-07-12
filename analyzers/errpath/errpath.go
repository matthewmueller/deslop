package errpath

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// New returns an analyzer that flags "if err == nil" patterns
// and unnecessary else blocks after returns.
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

		// Check for "if err == nil" (but allow "nil == err" as an intentional escape hatch).
		if isErrEqualsNil(ifStmt.Cond) {
			pass.Reportf(ifStmt.Pos(), "avoid \"if err == nil\"; invert to \"if err != nil { return ... }\" and keep the happy path on the left. If the nil path is genuinely shorter, use \"nil == err\" to suppress this warning.")
		}

		// Check for unnecessary else after a terminating if-body.
		if ifStmt.Else != nil && bodyTerminates(ifStmt.Body) {
			pass.Reportf(ifStmt.Else.Pos(), "unnecessary else after return/continue/break; remove the else and dedent the code")
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

	// Check "err == nil" (left is err, right is nil).
	if isErrIdent(binExpr.X) && isNilIdent(binExpr.Y) {
		return true
	}

	// "nil == err" is allowed as an escape hatch — don't flag it.
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

// bodyTerminates checks if a block ends with a return, continue, or break.
func bodyTerminates(body *ast.BlockStmt) bool {
	if len(body.List) == 0 {
		return false
	}
	last := body.List[len(body.List)-1]
	switch last.(type) {
	case *ast.ReturnStmt:
		return true
	case *ast.BranchStmt:
		return true
	}
	return false
}
