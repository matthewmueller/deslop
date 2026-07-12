package testcheck

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// NewIsCheck reports direct t.Fatal/t.Error calls and discouraged test imports.
func NewIsCheck() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "ischeck",
		Doc:      "reports direct t.Fatal/t.Error calls in tests; use github.com/matryer/is instead",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      runIsCheck,
	}
}

const isTemplate = `Use github.com/matryer/is for test assertions. The API:
  - is.Equal(got, want) asserts got == want (uses reflect.DeepEqual for structs/slices)
  - is.True(expr) asserts expr is true
  - is.NoErr(err) asserts err is nil

=== TEMPLATE START ===

package example_test

import (
    "testing"

    "github.com/matryer/is"
)

func TestCreate(t *testing.T) {
    is := is.New(t)
    result, err := Create("test")
    is.NoErr(err)
    is.Equal(result.Name, "test")
}

func TestCreateValidation(t *testing.T) {
    is := is.New(t)
    _, err := Create("")
    is.True(err != nil)
    is.True(strings.Contains(err.Error(), "name is required"))
}

// Use a helper for repeated setup:
func setup(t *testing.T) *Client {
    t.Helper()
    client, err := NewClient("test://localhost")
    if err != nil {
        t.Fatal(err) // t.Fatal is OK in setup helpers, not in assertions
    }
    return client
}

func TestClientGet(t *testing.T) {
    is := is.New(t)
    client := setup(t)
    resp, err := client.Get("/health")
    is.NoErr(err)
    is.Equal(resp.Status, 200)
}

=== TEMPLATE END ===`

func runIsCheck(pass *analysis.Pass) (any, error) {
	if !hasTestFiles(pass) {
		return nil, nil
	}

	flaggedImports := map[string]bool{
		"github.com/stretchr/testify":         true,
		"github.com/stretchr/testify/assert":  true,
		"github.com/stretchr/testify/require": true,
	}

	flaggedMethods := map[string]bool{
		"Fatal":  true,
		"Fatalf": true,
		"Error":  true,
		"Errorf": true,
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
		(*ast.ImportSpec)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		switch node := n.(type) {
		case *ast.ImportSpec:
			path := strings.Trim(node.Path.Value, `"`)
			if flaggedImports[path] {
				pass.Reportf(node.Pos(), "do not use %q; use github.com/matryer/is instead\n\n%s", path, isTemplate)
			}
		case *ast.CallExpr:
			sel, ok := node.Fun.(*ast.SelectorExpr)
			if !ok {
				return
			}
			if !flaggedMethods[sel.Sel.Name] {
				return
			}
			ident, ok := sel.X.(*ast.Ident)
			if !ok {
				return
			}
			if ident.Name != "t" && ident.Name != "tb" {
				return
			}
			pos := pass.Fset.Position(node.Pos())
			if !strings.HasSuffix(pos.Filename, "_test.go") {
				return
			}
			if isInHelper(pass, node) {
				return
			}
			pass.Reportf(node.Pos(), "avoid t.%s for assertions; use github.com/matryer/is instead\n\n%s", sel.Sel.Name, isTemplate)
		}
	})

	return nil, nil
}

// isInHelper checks if the call is inside a function that calls t.Helper().
func isInHelper(pass *analysis.Pass, call *ast.CallExpr) bool {
	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			if fn.Pos() <= call.Pos() && call.End() <= fn.End() {
				return hasHelperCall(fn.Body)
			}
		}
	}
	return false
}

func hasHelperCall(body *ast.BlockStmt) bool {
	for _, stmt := range body.List {
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}
		call, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			continue
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		if sel.Sel.Name == "Helper" {
			return true
		}
	}
	return false
}
