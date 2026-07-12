package noglobals

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "noglobals",
	Doc:  "reports package-level var declarations; use constructor injection instead",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		// Skip test files.
		name := pass.Fset.File(f.Pos()).Name()
		if strings.HasSuffix(name, "_test.go") {
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
					if isAllowed(ident.Name, vs) {
						continue
					}
					pass.Reportf(ident.Pos(), "avoid package-level var %q; use constructor injection instead", ident.Name)
				}
			}
		}
	}
	return nil, nil
}

func isAllowed(name string, spec *ast.ValueSpec) bool {
	// Allow Err* sentinels (e.g., var ErrNotFound = errors.New("..."))
	if strings.HasPrefix(name, "Err") {
		return true
	}
	// Allow blank identifier compile-time checks (var _ Interface = (*Type)(nil))
	if name == "_" {
		return true
	}
	return false
}
