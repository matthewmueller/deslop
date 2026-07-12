package complexity

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const defaultMaxLineLength = 120

// NewLineLength flags lines exceeding a character length threshold,
// ignoring lines where the length is due to a string literal.
func NewLineLength() *analysis.Analyzer {
	a := &analysis.Analyzer{
		Name: "linelength",
		Doc:  "reports lines that exceed the maximum character length",
	}
	var maxLength int
	a.Flags.IntVar(&maxLength, "max", defaultMaxLineLength, "maximum allowed line length in characters")
	a.Run = func(pass *analysis.Pass) (any, error) {
		return runLineLength(pass, maxLength)
	}
	return a
}

func runLineLength(pass *analysis.Pass, maxLength int) (any, error) {
	for _, f := range pass.Files {
		tokFile := pass.Fset.File(f.Pos())
		filename := tokFile.Name()
		if strings.Contains(filename, "/go-build/") {
			continue
		}

		stringLines := collectStringLines(pass.Fset, f)

		lineCount := tokFile.LineCount()
		for i := 1; i <= lineCount; i++ {
			lineStart := tokFile.LineStart(i)
			var lineEnd token.Pos
			if i < lineCount {
				lineEnd = tokFile.LineStart(i+1) - 1
			} else {
				lineEnd = token.Pos(tokFile.Base() + tokFile.Size())
			}
			length := int(lineEnd) - int(lineStart)
			if length <= maxLength {
				continue
			}
			if stringLines[i] {
				continue
			}
			msg := fmt.Sprintf("line is %d characters long (max %d); break it up", length, maxLength)
			pass.Reportf(lineStart, "%s", msg)
		}
	}
	return nil, nil
}

func collectStringLines(fset *token.FileSet, f *ast.File) map[int]bool {
	lines := map[int]bool{}
	ast.Inspect(f, func(n ast.Node) bool {
		lit, ok := n.(*ast.BasicLit)
		if !ok {
			return true
		}
		if lit.Kind != token.STRING {
			return true
		}
		pos := fset.Position(lit.Pos())
		end := fset.Position(lit.End())
		for line := pos.Line; line <= end.Line; line++ {
			lines[line] = true
		}
		return true
	})
	return lines
}
