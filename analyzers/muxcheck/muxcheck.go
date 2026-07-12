package muxcheck

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const template = `Use github.com/livebud/mux for HTTP routing.

The API works like this:
  - router := mux.New() creates a new router
  - router.Get(pattern, http.Handler) registers a GET route
  - router.Post(pattern, http.Handler) registers a POST route
  - router.Put(pattern, http.Handler) registers a PUT route
  - router.Patch(pattern, http.Handler) registers a PATCH route
  - router.Delete(pattern, http.Handler) registers a DELETE route
  - router.ServeHTTP(w, r) implements http.Handler
  - r.URL.Query().Get("param") extracts a path parameter inside a handler

Handlers must satisfy http.Handler. Use http.HandlerFunc to wrap a function.

Patterns use {param} syntax for path parameters:
  - "/users/{id}" matches /users/123, extract with r.URL.Query().Get("id")
  - "/{path*}" wildcard, matches remaining path segments
  - "/{major|[0-9]+}" regexp-constrained slot
  - "/fly/{from}-{to}" smart delimiters between slots

=== TEMPLATE START ===

package app

import (
    "net/http"

    "github.com/livebud/mux"
)

func NewRouter() http.Handler {
    router := mux.New()

    router.Get("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }))

    router.Get("/users/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id := r.URL.Query().Get("id")
        // use id...
        _ = id
    }))

    router.Post("/users", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // create user...
    }))

    // Middleware: mux.Use wraps func(next http.Handler) http.Handler
    router.Use(mux.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // before request...
            next.ServeHTTP(w, r)
            // after request...
        })
    }))

    return router
}

=== TEMPLATE END ===`

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "muxcheck",
		Doc:      "reports use of discouraged HTTP routers; use github.com/livebud/mux instead",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      run,
	}
}

func run(pass *analysis.Pass) (any, error) {
	flaggedImports := map[string]bool{
		"github.com/gorilla/mux":              true,
		"github.com/julienschmidt/httprouter": true,
		"github.com/go-chi/chi":               true,
		"github.com/go-chi/chi/v5":            true,
		"github.com/gin-gonic/gin":            true,
		"github.com/labstack/echo":            true,
		"github.com/labstack/echo/v4":         true,
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.ImportSpec)(nil),
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		switch node := n.(type) {
		case *ast.ImportSpec:
			path := strings.Trim(node.Path.Value, `"`)
			if flaggedImports[path] {
				pass.Reportf(node.Pos(), "do not use %q; use github.com/livebud/mux instead\n\n%s", path, template)
			}
		case *ast.CallExpr:
			checkMuxCall(pass, node)
		}
	})

	return nil, nil
}

func checkMuxCall(pass *analysis.Pass, node *ast.CallExpr) {
	sel, ok := node.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return
	}
	if sel.Sel.Name != "NewServeMux" {
		return
	}
	if !isPackageIdent(pass, ident, "net/http") {
		return
	}
	pass.Reportf(node.Pos(), "do not use http.NewServeMux; use github.com/livebud/mux instead\n\n%s", template)
}

func isPackageIdent(pass *analysis.Pass, ident *ast.Ident, pkgPath string) bool {
	obj := pass.TypesInfo.Uses[ident]
	if obj == nil {
		return false
	}
	pkgName, ok := obj.(*types.PkgName)
	if !ok {
		return false
	}
	return pkgName.Imported().Path() == pkgPath
}
