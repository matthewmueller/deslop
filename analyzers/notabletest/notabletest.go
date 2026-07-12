package notabletest

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "notabletest",
	Doc:      "reports table-driven tests; use top-level TestX functions with helpers instead",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	// Only check test files.
	hasTestFile := false
	for _, f := range pass.Files {
		name := pass.Fset.File(f.Pos()).Name()
		if strings.HasSuffix(name, "_test.go") {
			hasTestFile = true
			break
		}
	}
	if !hasTestFile {
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
		// Only check Test* functions.
		if !strings.HasPrefix(fn.Name.Name, "Test") {
			return
		}
		// Only check functions in test files.
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
		// Check if the range is over a variable that was assigned a []struct{} literal.
		if isTableLiteral(rangeStmt.X, body) {
			pass.Reportf(rangeStmt.Pos(), "avoid table-driven tests; use top-level TestX functions with helper functions instead")
		}
	}
}

// isTableLiteral checks if expr is (or refers to) a []struct{...} composite literal.
func isTableLiteral(expr ast.Expr, body *ast.BlockStmt) bool {
	switch x := expr.(type) {
	case *ast.CompositeLit:
		return isSliceOfStruct(x.Type)
	case *ast.Ident:
		// Look for the assignment in the same block.
		for _, stmt := range body.List {
			switch s := stmt.(type) {
			case *ast.AssignStmt:
				for i, lhs := range s.Lhs {
					ident, ok := lhs.(*ast.Ident)
					if !ok || ident.Name != x.Name {
						continue
					}
					if i < len(s.Rhs) {
						if lit, ok := s.Rhs[i].(*ast.CompositeLit); ok {
							return isSliceOfStruct(lit.Type)
						}
					}
				}
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
