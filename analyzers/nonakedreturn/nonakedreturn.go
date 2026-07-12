package nonakedreturn

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "nonakedreturn",
		Doc:  "reports bare return statements in functions with named return values",
		Run:  run,
	}
}

func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}
			if !hasNamedResults(fn) {
				return true
			}
			checkBody(pass, fn.Body)
			return true
		})
	}
	return nil, nil
}

func hasNamedResults(fn *ast.FuncDecl) bool {
	if fn.Type.Results == nil {
		return false
	}
	for _, field := range fn.Type.Results.List {
		if len(field.Names) > 0 {
			return true
		}
	}
	return false
}

func checkBody(pass *analysis.Pass, body *ast.BlockStmt) {
	if body == nil {
		return
	}
	ast.Inspect(body, func(n ast.Node) bool {
		// Skip nested function literals — they have their own returns.
		if _, ok := n.(*ast.FuncLit); ok {
			return false
		}
		ret, ok := n.(*ast.ReturnStmt)
		if !ok {
			return true
		}
		if len(ret.Results) == 0 {
			pass.Reportf(ret.Pos(), "avoid naked return; explicitly return values for readability")
		}
		return true
	})
}
