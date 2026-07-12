package contextfirst

import (
	"go/ast"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "contextfirst",
		Doc:  "reports exported functions/methods where context.Context is not the first parameter",
		Run:  run,
	}
}

func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		name := pass.Fset.File(f.Pos()).Name()
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Name == nil {
				continue
			}
			if !isExported(fn.Name.Name) {
				continue
			}
			checkFunc(pass, fn)
		}
	}
	return nil, nil
}

func checkFunc(pass *analysis.Pass, fn *ast.FuncDecl) {
	params := fn.Type.Params
	if params == nil || len(params.List) == 0 {
		return
	}

	// Check if first param is context.Context — if so, no issue.
	if isContextContext(params.List[0].Type) {
		return
	}

	// Check remaining params for context.Context.
	for i := 1; i < len(params.List); i++ {
		if isContextContext(params.List[i].Type) {
			pass.Reportf(fn.Pos(), "context.Context should be the first parameter of %s", fn.Name.Name)
			return
		}
	}
}

func isContextContext(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == "context" && sel.Sel.Name == "Context"
}

func isExported(name string) bool {
	if name == "" {
		return false
	}
	return unicode.IsUpper(rune(name[0]))
}
