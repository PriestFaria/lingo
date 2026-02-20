package log

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)
type LogContext struct {
	Pass *analysis.Pass
	CallExpr *ast.CallExpr
	Parts []string 
	FullText string
	Suggestions []Suggestion
}