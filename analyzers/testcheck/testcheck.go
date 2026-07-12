package testcheck

import "golang.org/x/tools/go/analysis"

// New returns all testcheck analyzers.
func New() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		NewNoTable(),
		NewIsCheck(),
		NewNaming(),
	}
}
