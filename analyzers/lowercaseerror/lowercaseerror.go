package lowercaseerror

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "lowercaseerror",
		Doc:  "reports error strings that start with uppercase",
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
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			if !isErrorfCall(call) {
				return true
			}
			if len(call.Args) == 0 {
				return true
			}
			lit, ok := call.Args[0].(*ast.BasicLit)
			if !ok {
				return true
			}
			if lit.Kind != token.STRING {
				return true
			}
			val := lit.Value
			if len(val) < 2 {
				return true
			}
			// Remove opening quote
			content := val[1:]
			if len(content) == 0 {
				return true
			}
			firstRune := rune(content[0])
			if unicode.IsUpper(firstRune) {
				msg := fmt.Sprintf("error strings should start with lowercase; got %s", val)
				pass.Reportf(lit.Pos(), "%s", msg)
			}
			return true
		})
	}
	return nil, nil
}

func isErrorfCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == "fmt" && sel.Sel.Name == "Errorf"
}
