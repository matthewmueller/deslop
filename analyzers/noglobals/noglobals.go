package noglobals

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "noglobals",
		Doc:  "reports package-level var declarations; use constructor injection instead",
		Run:  run,
	}
}

func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		name := pass.Fset.File(f.Pos()).Name()
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		// Skip generated files (e.g., Go's test main from the build cache).
		if !strings.HasSuffix(name, ".go") || strings.Contains(name, "/go-build/") {
			continue
		}

		for _, decl := range f.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.VAR {
				continue
			}
			for _, spec := range gen.Specs {
				vs := spec.(*ast.ValueSpec)
				for _, ident := range vs.Names {
					if isAllowed(ident.Name) {
						continue
					}
					pass.Reportf(ident.Pos(), "avoid package-level var %q; use constructor injection instead", ident.Name)
				}
			}
		}
	}
	return nil, nil
}

func isAllowed(name string) bool {
	// Allow unexported (private) vars.
	if len(name) > 0 && name[0] >= 'a' && name[0] <= 'z' {
		return true
	}
	// Allow Err* sentinels.
	if strings.HasPrefix(name, "Err") {
		return true
	}
	// Allow blank identifier compile-time checks.
	if name == "_" {
		return true
	}
	return false
}
