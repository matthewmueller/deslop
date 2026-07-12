package envcheck

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const template = `package env

import (
    env11 "github.com/caarlos0/env/v11"
)

// Env environment
type Env struct {
    %s string ` + "`" + `env:"%s,required"` + "`" + `
}

// Load the environment
func Load() (*Env, error) {
    env := new(Env)
    if err := env11.Parse(env); err != nil {
        return nil, err
    }
    return env, nil
}`

func New() *analysis.Analyzer {
	a := &analysis.Analyzer{
		Name:     "envcheck",
		Doc:      "reports direct os.Getenv/os.LookupEnv calls; use internal/env instead",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
	var includeTests bool
	a.Flags.BoolVar(&includeTests, "tests", false, "check test files too")
	a.Run = func(pass *analysis.Pass) (any, error) {
		return run(pass, includeTests)
	}
	return a
}

func run(pass *analysis.Pass, includeTests bool) (any, error) {
	if strings.HasSuffix(pass.Pkg.Path(), "internal/env") {
		return nil, nil
	}
	if !includeTests && isTestPackage(pass) {
		return nil, nil
	}

	flagged := map[string]bool{
		"Getenv":    true,
		"LookupEnv": true,
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}
		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			return
		}
		if ident.Name != "os" || !flagged[sel.Sel.Name] {
			return
		}

		envVar := extractStringArg(call)
		msg := fmt.Sprintf(
			"do not use os.%s directly; add %s to the Env struct in internal/env/env.go and load it via env.Load()\n\n"+
				"If internal/env/env.go does not exist, create it with this template:\n\n%s",
			sel.Sel.Name,
			envVar,
			fmt.Sprintf(template, envVar, envVar),
		)
		pass.Reportf(call.Pos(), "%s", msg)
	})

	return nil, nil
}

func isTestPackage(pass *analysis.Pass) bool {
	for _, f := range pass.Files {
		name := pass.Fset.File(f.Pos()).Name()
		if strings.HasSuffix(name, "_test.go") {
			return true
		}
	}
	return false
}

func extractStringArg(call *ast.CallExpr) string {
	if len(call.Args) == 0 {
		return "UNKNOWN"
	}
	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok {
		return "UNKNOWN"
	}
	return strings.Trim(lit.Value, `"`)
}
