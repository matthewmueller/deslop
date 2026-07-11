package clicheck

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "clicheck",
	Doc:      "reports direct os.Args access and flag parsing outside internal/cli",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

// flaggedImports contains package paths that indicate direct CLI parsing.
var flaggedImports = map[string]bool{
	"flag": true,
}

const template = `If internal/cli does not exist, create internal/cli/cli.go with this template.

The CLI struct is the central object. Each subcommand gets:
  1. An input struct with exported fields
  2. A method on *CLI that takes context + the input struct
  3. A block in Parse() that wires flags/args to the input struct

Use github.com/livebud/cli for parsing. The API works like this:
  - cli.New(name, description) creates the root command
  - cmd.Command(name, description) creates a subcommand
  - cmd.Flag(name, description).String(&target) binds a flag
  - cmd.Flag(name, description).Enum(&target, ...values).Default(value) binds an enum flag
  - cmd.Flag(name, description).Bool(&target).Default(value) binds a bool flag
  - cmd.Flag(name, description).Int(&target).Default(value) binds an int flag
  - cmd.Flag(name, description).Optional().String(&target) binds an optional flag (target is **string, nil when unset)
  - cmd.Arg(name, description).String(&target) binds a positional arg
  - cmd.Arg(name, description).Strings(&target) binds variadic trailing args
  - cmd.Args(name, description).Strings(&target) binds all remaining args as a string slice
  - cmd.Run(func(ctx context.Context) error { ... }) sets the handler

=== TEMPLATE START ===

package cli

import (
    "context"
    "fmt"
    "io"
    "log/slog"
    "os"

    "github.com/livebud/cli"
    "github.com/matthewmueller/logs"
)

func Run() int {
    cli := Default()
    ctx := context.Background()
    err := cli.Parse(ctx, os.Args[1:]...)
    if err != nil {
        logs.ErrorContext(ctx, err.Error())
        return 1
    }
    return 0
}

func Default() *CLI {
    return &CLI{
        Stdout:   os.Stdout,
        Stderr:   os.Stderr,
        Stdin:    os.Stdin,
        Env:      os.Environ(),
        Dir:      ".",
        logLevel: "info",
    }
}

type CLI struct {
    Stdout io.Writer
    Stderr io.Writer
    Stdin  io.Reader
    Env    []string
    Dir    string

    logLevel string
}

func (c *CLI) log() (*slog.Logger, error) {
    level, err := logs.ParseLevel(c.logLevel)
    if err != nil {
        return nil, fmt.Errorf("cli: parsing log level: %w", err)
    }
    log := logs.New(logs.Filter(level, logs.Console(c.Stdout)))
    return log, nil
}

func (c *CLI) Parse(ctx context.Context, args ...string) error {
    cli := cli.New("app", "app description")
    cli.Flag("log", "log level").Enum(&c.logLevel, "debug", "info", "warn", "error").Default("info")

    { // create <name>
        in := &Create{}
        cmd := cli.Command("create", "create a new resource")
        cmd.Arg("name", "name of the resource").String(&in.Name)
        cmd.Run(func(ctx context.Context) error {
            return c.Create(ctx, in)
        })
    }

    { // upload [--tag=<tag>] [--force] <from> <to>
        in := &Upload{}
        cmd := cli.Command("upload", "upload a file")
        cmd.Flag("tag", "tag to apply").Optional().String(&in.Tag)
        cmd.Flag("force", "overwrite existing").Bool(&in.Force).Default(false)
        cmd.Arg("from", "source path").String(&in.From)
        cmd.Arg("to", "destination path").String(&in.To)
        cmd.Run(func(ctx context.Context) error {
            return c.Upload(ctx, in)
        })
    }

    // Nested subcommands:
    cache := cli.Command("cache", "cache operations")

    { // cache prune [--older-than=<duration>] <path>
        in := &CachePrune{}
        cmd := cache.Command("prune", "prune the cache")
        cmd.Flag("older-than", "prune entries older than duration").Optional().String(&in.OlderThan)
        cmd.Arg("path", "cache path").String(&in.Path)
        cmd.Run(func(ctx context.Context) error {
            return c.CachePrune(ctx, in)
        })
    }

    return cli.Parse(ctx, args...)
}

// --- Input structs ---

type Create struct {
    Name string
}

func (c *CLI) Create(ctx context.Context, in *Create) error {
    return nil
}

type Upload struct {
    Tag   string
    Force bool
    From  string
    To    string
}

func (c *CLI) Upload(ctx context.Context, in *Upload) error {
    return nil
}

type CachePrune struct {
    OlderThan string
    Path      string
}

func (c *CLI) CachePrune(ctx context.Context, in *CachePrune) error {
    return nil
}

=== TEMPLATE END ===`

func run(pass *analysis.Pass) (any, error) {
	if isAllowed(pass.Pkg.Path()) {
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.SelectorExpr)(nil),
		(*ast.ImportSpec)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		switch node := n.(type) {
		case *ast.SelectorExpr:
			ident, ok := node.X.(*ast.Ident)
			if !ok {
				return
			}
			if ident.Name == "os" && node.Sel.Name == "Args" {
				pass.Reportf(node.Pos(), "do not access os.Args directly; use internal/cli instead\n\n%s", template)
			}
		case *ast.ImportSpec:
			path := strings.Trim(node.Path.Value, `"`)
			if flaggedImports[path] {
				pass.Reportf(node.Pos(), "do not use the %q package directly; use internal/cli instead\n\n%s", path, template)
			}
		}
	})

	return nil, nil
}

func isAllowed(pkgPath string) bool {
	return strings.HasSuffix(pkgPath, "internal/cli") ||
		strings.Contains(pkgPath, "internal/cli/")
}
