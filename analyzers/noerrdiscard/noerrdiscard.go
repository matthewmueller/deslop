package noerrdiscard

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "noerrdiscard",
		Doc:  "reports discarded errors via _ assignment from function calls",
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
		reported := map[*ast.AssignStmt]bool{}
		ast.Inspect(f, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.IfStmt:
				checkIfInit(pass, node, reported)
			case *ast.AssignStmt:
				if !reported[node] {
					checkAssign(pass, node)
				}
			}
			return true
		})
	}
	return nil, nil
}

func checkAssign(pass *analysis.Pass, assign *ast.AssignStmt) {
	if !isDiscardingError(pass, assign) {
		return
	}
	pass.Reportf(assign.Pos(), "do not discard errors; handle the error or propagate it")
}

func checkIfInit(pass *analysis.Pass, ifStmt *ast.IfStmt, reported map[*ast.AssignStmt]bool) {
	assign, ok := ifStmt.Init.(*ast.AssignStmt)
	if !ok {
		return
	}
	reported[assign] = true
	if !isDiscardingError(pass, assign) {
		return
	}
	pass.Reportf(assign.Pos(), "do not discard errors; handle the error or propagate it")
}

func isDiscardingError(pass *analysis.Pass, assign *ast.AssignStmt) bool {
	if assign.Tok != token.ASSIGN && assign.Tok != token.DEFINE {
		return false
	}
	if len(assign.Rhs) != 1 {
		return false
	}
	if !isCallExpr(assign.Rhs[0]) {
		return false
	}
	// Last LHS must be blank identifier.
	last := assign.Lhs[len(assign.Lhs)-1]
	ident, ok := last.(*ast.Ident)
	if !ok || ident.Name != "_" {
		return false
	}
	// Check if the discarded value is actually an error type.
	return lastReturnIsError(pass, assign.Rhs[0])
}

func lastReturnIsError(pass *analysis.Pass, expr ast.Expr) bool {
	typ := pass.TypesInfo.TypeOf(expr)
	if typ == nil {
		return false
	}
	// Single return value.
	if isErrorType(typ) {
		return true
	}
	// Multi-return: check the last element of the tuple.
	tuple, ok := typ.(*types.Tuple)
	if !ok {
		return false
	}
	if tuple.Len() == 0 {
		return false
	}
	lastType := tuple.At(tuple.Len() - 1).Type()
	return isErrorType(lastType)
}

func isErrorType(t types.Type) bool {
	return types.Implements(t, errorInterface())
}

func errorInterface() *types.Interface {
	return types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
}

func isCallExpr(expr ast.Expr) bool {
	_, ok := expr.(*ast.CallExpr)
	return ok
}
