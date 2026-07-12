package nolog

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "nolog",
		Doc:      "reports use of stdlib log package; use log/slog with structured fields instead",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      run,
	}
}

const template = `Use log/slog with structured fields and package-prefixed messages.
Initialize loggers with github.com/matthewmueller/logs.

=== TEMPLATE START ===

import (
    "log/slog"

    "github.com/matthewmueller/logs"
)

// Initialize the logger (typically in main or CLI setup).
// logs.Default() and logs.New() both return a *slog.Logger.

// Go with the default logger:
log := logs.Default()

// Or configure:
log = logs.New(
    logs.Multi(
        logs.Filter(slog.LevelInfo, logs.Console(os.Stderr)),
        slog.NewJSONHandler(os.Stderr, nil),
    ),
)

// Accept *slog.Logger via constructor injection:
type Server struct {
    log *slog.Logger
}

func NewServer(log *slog.Logger) *Server {
    return &Server{log: log}
}

// Use structured fields and prefix messages with package context:
func (s *Server) Handle(ctx context.Context, r *http.Request) {
    s.log.InfoContext(ctx, "server: handling request", "method", r.Method, "path", r.URL.Path)

    if err != nil {
        s.log.ErrorContext(ctx, "server: failed to process", "error", err, "path", r.URL.Path)
    }
}

=== TEMPLATE END ===`

func run(pass *analysis.Pass) (any, error) {
	flaggedCalls := map[string]bool{
		"Print":   true,
		"Printf":  true,
		"Println": true,
		"Fatal":   true,
		"Fatalf":  true,
		"Fatalln": true,
		"Panic":   true,
		"Panicf":  true,
		"Panicln": true,
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
			if path == "log" {
				pos := pass.Fset.Position(node.Pos())
				if strings.HasSuffix(pos.Filename, "_test.go") {
					return
				}
				pass.Reportf(node.Pos(), "do not use the stdlib \"log\" package; use log/slog instead\n\n%s", template)
			}
		case *ast.CallExpr:
			sel, ok := node.Fun.(*ast.SelectorExpr)
			if !ok {
				return
			}
			ident, ok := sel.X.(*ast.Ident)
			if !ok {
				return
			}
			if ident.Name == "log" && flaggedCalls[sel.Sel.Name] {
				pos := pass.Fset.Position(node.Pos())
				if strings.HasSuffix(pos.Filename, "_test.go") {
					return
				}
				pass.Reportf(node.Pos(), "do not use log.%s; use log/slog with structured fields instead\n\n%s", sel.Sel.Name, template)
			}
		}
	})

	return nil, nil
}
