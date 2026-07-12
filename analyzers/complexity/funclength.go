package complexity

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const defaultMaxFuncLength = 500

// NewFuncLength flags functions exceeding a line count threshold.
func NewFuncLength() *analysis.Analyzer {
	a := &analysis.Analyzer{
		Name: "funclength",
		Doc:  "reports functions that are too long",
	}
	var maxLength int
	a.Flags.IntVar(&maxLength, "max", defaultMaxFuncLength, "maximum allowed function length in lines")
	a.Run = func(pass *analysis.Pass) (any, error) {
		return runFuncLength(pass, maxLength)
	}
	return a
}

func runFuncLength(pass *analysis.Pass, maxLength int) (any, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				return true
			}
			start := pass.Fset.Position(fn.Pos()).Line
			end := pass.Fset.Position(fn.End()).Line
			length := end - start + 1
			if length > maxLength {
				msg := fmt.Sprintf("function %s is %d lines long (max %d); break it into smaller functions", fn.Name.Name, length, maxLength)
				pass.Reportf(fn.Name.Pos(), "%s", msg)
			}
			return false
		})
	}
	return nil, nil
}
