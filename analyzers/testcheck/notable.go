package testcheck

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// NewNoTable reports table-driven tests.
func NewNoTable() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "notabletest",
		Doc:      "reports table-driven tests; inline assertions or break into separate TestX functions",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      runNoTable,
	}
}

func runNoTable(pass *analysis.Pass) (any, error) {
	if !hasTestFiles(pass) {
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		fn := n.(*ast.FuncDecl)
		if fn.Body == nil {
			return
		}
		if !strings.HasPrefix(fn.Name.Name, "Test") {
			return
		}
		pos := pass.Fset.Position(fn.Pos())
		if !strings.HasSuffix(pos.Filename, "_test.go") {
			return
		}
		checkForTableTest(pass, fn.Body)
	})

	return nil, nil
}

func checkForTableTest(pass *analysis.Pass, body *ast.BlockStmt) {
	for _, stmt := range body.List {
		rangeStmt, ok := stmt.(*ast.RangeStmt)
		if !ok {
			continue
		}
		if isTableLiteral(rangeStmt.X, body) {
			pass.Reportf(rangeStmt.Pos(), "avoid table-driven tests; instead either inline the assertions in a single test function or break into separate TestX functions if the cases are meaningfully different")
		}
	}
}

func isTableLiteral(expr ast.Expr, body *ast.BlockStmt) bool {
	switch x := expr.(type) {
	case *ast.CompositeLit:
		return isSliceOfStruct(x.Type)
	case *ast.Ident:
		return identRefersToStructSlice(x.Name, body)
	}
	return false
}

func identRefersToStructSlice(name string, body *ast.BlockStmt) bool {
	for _, stmt := range body.List {
		assign, ok := stmt.(*ast.AssignStmt)
		if !ok {
			continue
		}
		for i, lhs := range assign.Lhs {
			ident, ok := lhs.(*ast.Ident)
			if !ok || ident.Name != name {
				continue
			}
			if i >= len(assign.Rhs) {
				continue
			}
			lit, ok := assign.Rhs[i].(*ast.CompositeLit)
			if !ok {
				continue
			}
			if isSliceOfStruct(lit.Type) {
				return true
			}
		}
	}
	return false
}

func isSliceOfStruct(expr ast.Expr) bool {
	arr, ok := expr.(*ast.ArrayType)
	if !ok {
		return false
	}
	_, ok = arr.Elt.(*ast.StructType)
	return ok
}
