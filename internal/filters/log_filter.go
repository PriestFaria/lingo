package filters

import (
	"github.com/PriestFaria/lingo/internal/analyzer/log"
)

// LogFilter is the interface implemented by all log message filters.
// Apply inspects the provided LogContext and returns a slice of issues found.
// An empty slice means no violations.
type LogFilter interface {
	Apply(context *log.LogContext) []FilterIssue
}

