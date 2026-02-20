package log

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

type LogPart struct {
	Value string
	IsLiteral bool
	Pos token.Pos
	End token.Pos

}
type LogContext struct {
	Pass *analysis.Pass
	CallExpr *ast.CallExpr
	Parts []LogPart 
	FullText string
}