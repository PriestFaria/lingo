package filters

import "lingo/internal/analyzer/log"

type SecurityFilter struct{}

func (f *SecurityFilter) Apply(context *log.LogContext) []FilterIssue {
	return []FilterIssue{}
}