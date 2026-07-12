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
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if fn.Recv == nil {
				continue
			}
			methodName := fn.Name.Name
			if !strings.HasPrefix(methodName, "Get") {
				continue
			}
			suffix := strings.TrimPrefix(methodName, "Get")
			if len(suffix) == 0 {
				continue
			}
			firstRune := rune(suffix[0])
			if !unicode.IsUpper(firstRune) {
				continue
			}
			msg := fmt.Sprintf("method %s should be renamed to %s; avoid Get prefixes", methodName, suffix)
			pass.Reportf(fn.Name.Pos(), "%s", msg)
		}
	}
	return nil, nil
}
