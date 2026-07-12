package complexity

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const defaultMaxNameLength = 30

// NewNaming flags functions with names exceeding a length threshold.
func NewNaming() *analysis.Analyzer {
	a := &analysis.Analyzer{
		Name: "funcnaming",
		Doc:  "reports functions with names that are too long",
	}
	var maxNameLength int
	a.Flags.IntVar(&maxNameLength, "max", defaultMaxNameLength, "maximum allowed function name length")
	a.Run = func(pass *analysis.Pass) (any, error) {
		return runNaming(pass, maxNameLength)
	}
	return a
}

func runNaming(pass *analysis.Pass, maxNameLength int) (any, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}
			name := fn.Name.Name
			if len(name) > maxNameLength {
				msg := fmt.Sprintf("function name %q is %d characters (max %d); use a shorter name", name, len(name), maxNameLength)
				pass.Reportf(fn.Name.Pos(), "%s", msg)
			}
			return false
		})
	}
	return nil, nil
}
