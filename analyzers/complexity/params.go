package complexity

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const defaultMaxParams = 4

// NewParams flags functions with too many parameters.
func NewParams() *analysis.Analyzer {
	a := &analysis.Analyzer{
		Name: "params",
		Doc:  "reports functions with too many parameters; use an options struct instead",
	}
	var maxParams int
	a.Flags.IntVar(&maxParams, "max", defaultMaxParams, "maximum allowed parameter count")
	a.Run = func(pass *analysis.Pass) (any, error) {
		return runParams(pass, maxParams)
	}
	return a
}

func runParams(pass *analysis.Pass, maxParams int) (any, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok || fn.Type.Params == nil {
				return true
			}
			count := 0
			for _, field := range fn.Type.Params.List {
				names := len(field.Names)
				if names == 0 {
					names = 1
				}
				count += names
			}
			if count > maxParams {
				msg := fmt.Sprintf("function %s has %d parameters (max %d); use an options struct instead", fn.Name.Name, count, maxParams)
				pass.Reportf(fn.Name.Pos(), "%s", msg)
			}
			return false
		})
	}
	return nil, nil
}
