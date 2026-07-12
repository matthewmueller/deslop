package returncheck

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "returncheck",
	Doc:  "reports New*/Load* constructors that return interfaces; return concrete types instead",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		name := pass.Fset.File(f.Pos()).Name()
		if strings.HasSuffix(name, "_test.go") {
			continue
		}

		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Type.Results == nil {
				continue
			}
			if !isConstructor(fn.Name.Name) {
				continue
			}
			// Check return types for interfaces.
			for _, result := range fn.Type.Results.List {
				if isInterfaceReturn(result.Type) {
					pass.Reportf(fn.Name.Pos(), "constructor %s should return a concrete type, not an interface", fn.Name.Name)
					break
				}
			}
		}
	}
	return nil, nil
}

func isConstructor(name string) bool {
	return strings.HasPrefix(name, "New") || strings.HasPrefix(name, "Load")
}

func isInterfaceReturn(expr ast.Expr) bool {
	switch t := expr.(type) {
	case *ast.InterfaceType:
		// Bare `interface{}` or `any`
		return true
	case *ast.Ident:
		// Named interface (heuristic: not a pointer, not a builtin like error)
		// We can't fully resolve types without type info, but we skip "error"
		// since returning error is standard.
		if t.Name == "error" {
			return false
		}
		if t.Obj != nil && t.Obj.Decl != nil {
			ts, ok := t.Obj.Decl.(*ast.TypeSpec)
			if ok {
				_, isIface := ts.Type.(*ast.InterfaceType)
				return isIface
			}
		}
	}
	return false
}
