package filters

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

type LogContext struct {
	Pass *analysis.Pass
	CallExpr *ast.CallExpr
	Parts []string 
	FullText string
}


type LogFilter interface {
	Apply(context *LogContext) bool
}