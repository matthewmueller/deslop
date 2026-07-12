package nocommentedcode

import (
	"go/ast"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "nocommentedcode",
		Doc:  "reports comments that look like commented-out code",
		Run:  run,
	}
}

// Patterns that indicate a non-code comment — skip these.
var skipPrefixes = []string{
	"TODO",
	"FIXME",
	"NOTE",
	"HACK",
	"XXX",
	"BUG",
	"nolint",
	"go:",
	"#",
}

// Go keywords that suggest commented-out code when followed by a space.
var codeKeywords = []string{
	"if ",
	"for ",
	"return ",
	"var ",
	"const ",
	"switch ",
	"select ",
	"defer ",
	"go ",
	"range ",
}

func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		checkComments(pass, f)
	}
	return nil, nil
}

func checkComments(pass *analysis.Pass, f *ast.File) {
	docComments := collectDocComments(f)
	for _, cg := range f.Comments {
		if docComments[cg] {
			continue
		}
		checkCommentGroup(pass, cg)
	}
}

func checkCommentGroup(pass *analysis.Pass, cg *ast.CommentGroup) {
	for _, c := range cg.List {
		if !strings.HasPrefix(c.Text, "//") {
			continue
		}
		text := strings.TrimSpace(strings.TrimPrefix(c.Text, "//"))
		if text == "" {
			continue
		}
		if isSkipPrefix(text) {
			continue
		}
		if looksLikeCode(text) {
			pass.Reportf(c.Pos(), "commented-out code should be deleted, not commented; use version control to recover old code")
		}
	}
}

func collectDocComments(f *ast.File) map[*ast.CommentGroup]bool {
	docs := map[*ast.CommentGroup]bool{}
	for _, decl := range f.Decls {
		doc := declDoc(decl)
		if doc != nil {
			docs[doc] = true
		}
	}
	return docs
}

func declDoc(decl ast.Decl) *ast.CommentGroup {
	switch d := decl.(type) {
	case *ast.GenDecl:
		return d.Doc
	case *ast.FuncDecl:
		return d.Doc
	}
	return nil
}

func isSkipPrefix(text string) bool {
	for _, p := range skipPrefixes {
		if strings.HasPrefix(text, p) {
			return true
		}
	}
	return false
}

func looksLikeCode(text string) bool {
	if looksLikeAssignment(text) {
		return true
	}
	if startsWithKeyword(text) {
		return true
	}
	if startsWithIdentDot(text) {
		return true
	}
	return false
}

func looksLikeAssignment(text string) bool {
	idx := strings.Index(text, ":=")
	if idx < 1 {
		return false
	}
	before := strings.TrimSpace(text[:idx])
	if before == "" {
		return false
	}
	for _, r := range before {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != ',' && r != ' ' {
			return false
		}
	}
	return true
}

func startsWithKeyword(text string) bool {
	for _, kw := range codeKeywords {
		if strings.HasPrefix(text, kw) {
			return true
		}
	}
	return false
}

func startsWithIdentDot(text string) bool {
	// Must start with a lowercase letter (identifier).
	if len(text) == 0 || !unicode.IsLower(rune(text[0])) {
		return false
	}
	// Find the dot.
	dotIdx := strings.IndexByte(text, '.')
	if dotIdx < 1 {
		return false
	}
	// Everything before dot should look like an identifier.
	ident := text[:dotIdx]
	for _, r := range ident {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	// After the dot, must look like a method call: ident.Method(
	rest := text[dotIdx+1:]
	parenIdx := strings.IndexByte(rest, '(')
	if parenIdx < 1 {
		return false
	}
	methodName := rest[:parenIdx]
	for _, r := range methodName {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}
