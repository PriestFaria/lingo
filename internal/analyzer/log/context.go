package log

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// LogPart represents a single segment of a log message argument.
// A message formed by concatenation is split into multiple parts.
type LogPart struct {
	// Value holds the literal text or the identifier name.
	Value string
	// IsLiteral is true for string literals and false for variables/expressions.
	IsLiteral bool
	Pos       token.Pos
	End       token.Pos
}
// LogContext carries the data a LogFilter needs to inspect a single log call.
type LogContext struct {
	Pass     *analysis.Pass
	CallExpr *ast.CallExpr
	Parts    []LogPart
	// FullText is the concatenation of all Part values for convenience.
	FullText string
}