package filters

import (
	"go/token"
	"lingo/internal/analyzer/log"
)

type Suggestion struct {
	Message string
	Replace string
	Start token.Pos
	End token.Pos
}



type LogFilter interface {
	Apply(context *log.LogContext) 
}