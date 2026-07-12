package noinit

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "noinit",
		Doc:  "reports init() function declarations",
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
			if fn.Name.Name == "init" && fn.Recv == nil {
				pass.Reportf(fn.Pos(), "avoid init() functions; use explicit initialization in constructors or main instead")
			}
		}
	}
	return nil, nil
}
