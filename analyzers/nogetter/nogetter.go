package nogetter

import (
	"fmt"
	"go/ast"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "nogetter",
		Doc:  "reports Get-prefixed method names",
		Run:  run,
	}
}

func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		name := pass.Fset.File(f.Pos()).Name()
		if strings.Contains(name, "/go-build/") {
			continue
		}
		for _, decl := range f.Decls {
			checkDecl(pass, decl)
		}
	}
	return nil, nil
}

func checkDecl(pass *analysis.Pass, decl ast.Decl) {
	fn, ok := decl.(*ast.FuncDecl)
	if !ok || fn.Recv == nil {
		return
	}
	methodName := fn.Name.Name
	if !strings.HasPrefix(methodName, "Get") {
		return
	}
	suffix := strings.TrimPrefix(methodName, "Get")
	if len(suffix) == 0 || !unicode.IsUpper(rune(suffix[0])) {
		return
	}
	if hasImplementsComment(fn) {
		return
	}
	msg := fmt.Sprintf(
		"method %s should be renamed to %s; avoid Get prefixes. "+
			"If this method implements an interface, add a doc comment that includes \"implements\" (e.g. // %s implements Interface).",
		methodName, suffix, methodName,
	)
	pass.Reportf(fn.Name.Pos(), "%s", msg)
}

func hasImplementsComment(fn *ast.FuncDecl) bool {
	if fn.Doc == nil {
		return false
	}
	for _, c := range fn.Doc.List {
		if strings.Contains(c.Text, "implements") || strings.Contains(c.Text, "Implements") {
			return true
		}
	}
	return false
}
