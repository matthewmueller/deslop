package complexity

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const defaultMaxNesting = 3

// NewNesting flags functions with nesting depth exceeding the threshold.
func NewNesting() *analysis.Analyzer {
	a := &analysis.Analyzer{
		Name: "nesting",
		Doc:  "reports functions with excessive nesting depth",
	}
	var maxNesting int
	a.Flags.IntVar(&maxNesting, "max", defaultMaxNesting, "maximum allowed nesting depth")
	a.Run = func(pass *analysis.Pass) (any, error) {
		return runNesting(pass, maxNesting)
	}
	return a
}

func runNesting(pass *analysis.Pass, maxNesting int) (any, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				return true
			}
			depth := maxDepth(fn.Body, 0)
			if depth > maxNesting {
				msg := fmt.Sprintf("function %s has nesting depth %d (max %d); simplify with early returns or extract helpers", fn.Name.Name, depth, maxNesting)
				pass.Reportf(fn.Name.Pos(), "%s", msg)
			}
			return false
		})
	}
	return nil, nil
}

func maxDepth(block *ast.BlockStmt, current int) int {
	max := current
	for _, stmt := range block.List {
		d := stmtDepth(stmt, current)
		if d > max {
			max = d
		}
	}
	return max
}

func stmtDepth(stmt ast.Stmt, current int) int {
	switch s := stmt.(type) {
	case *ast.IfStmt:
		return ifDepth(s, current)
	case *ast.ForStmt:
		return maxDepth(s.Body, current+1)
	case *ast.RangeStmt:
		return maxDepth(s.Body, current+1)
	case *ast.SwitchStmt:
		return maxDepth(s.Body, current+1)
	case *ast.TypeSwitchStmt:
		return maxDepth(s.Body, current+1)
	case *ast.SelectStmt:
		return maxDepth(s.Body, current+1)
	case *ast.CaseClause:
		return clauseDepth(s.Body, current)
	case *ast.CommClause:
		return clauseDepth(s.Body, current)
	case *ast.BlockStmt:
		return maxDepth(s, current+1)
	}
	return current
}

func ifDepth(s *ast.IfStmt, current int) int {
	max := maxDepth(s.Body, current+1)
	if s.Else == nil {
		return max
	}
	d := elseDepth(s.Else, current+1)
	if d > max {
		return d
	}
	return max
}

func elseDepth(stmt ast.Stmt, current int) int {
	switch s := stmt.(type) {
	case *ast.BlockStmt:
		return maxDepth(s, current)
	case *ast.IfStmt:
		return ifDepth(s, current)
	}
	return current
}

func clauseDepth(body []ast.Stmt, current int) int {
	max := current + 1
	for _, st := range body {
		d := stmtDepth(st, current+1)
		if d > max {
			max = d
		}
	}
	return max
}
