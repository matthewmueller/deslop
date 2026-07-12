package nostutter

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "nostutter",
		Doc:  "reports type names that stutter with package name",
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
		checkFile(pass, f)
	}
	return nil, nil
}

func checkFile(pass *analysis.Pass, f *ast.File) {
	pkgName := pass.Pkg.Name()
	prefix := strings.ToLower(pkgName)
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gen.Specs {
			checkSpec(pass, spec, pkgName, prefix)
		}
	}
}

func checkSpec(pass *analysis.Pass, spec ast.Spec, pkgName, prefix string) {
	ts, ok := spec.(*ast.TypeSpec)
	if !ok {
		return
	}
	typeName := ts.Name.Name
	if !ast.IsExported(typeName) {
		return
	}
	lower := strings.ToLower(typeName)
	if !strings.HasPrefix(lower, prefix) {
		return
	}
	suffix := typeName[len(pkgName):]
	if len(suffix) == 0 {
		return
	}
	msg := fmt.Sprintf(
		"type %s stutters with package name %s; consider renaming to %s",
		typeName, pkgName, suffix,
	)
	pass.Reportf(ts.Name.Pos(), "%s", msg)
}
