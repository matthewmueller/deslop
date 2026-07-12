package complexity

import "golang.org/x/tools/go/analysis"

// New returns all complexity analyzers.
func New() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		NewNesting(),
		NewParams(),
		NewNaming(),
		NewLineLength(),
		NewFuncLength(),
	}
}
