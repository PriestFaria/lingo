package filters

import (
	"lingo/internal/analyzer/log"
)


type LogFilter interface {
	Apply(context *log.LogContext) []FilterIssue 
}

