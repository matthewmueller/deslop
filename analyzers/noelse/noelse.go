package noelse

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "noelse",
		Doc:  "reports else blocks; prefer early returns or extract a function",
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
		ast.Inspect(f, func(n ast.Node) bool {
			ifStmt, ok := n.(*ast.IfStmt)
			if !ok {
				return true
			}
			if ifStmt.Else == nil {
				return true
			}
			// Allow "else if" chains that terminate (return/break/continue).
			if elseIf, ok := ifStmt.Else.(*ast.IfStmt); ok {
				if bodyTerminates(elseIf.Body) {
					return true
				}
			}
			pass.Reportf(ifStmt.Else.Pos(), "avoid else blocks; use an early return or extract a separate function instead")
			return true
		})
	}
	return nil, nil
}

func bodyTerminates(body *ast.BlockStmt) bool {
	if len(body.List) == 0 {
		return false
	}
	last := body.List[len(body.List)-1]
	switch last.(type) {
	case *ast.ReturnStmt, *ast.BranchStmt:
		return true
	}
	return false
}
