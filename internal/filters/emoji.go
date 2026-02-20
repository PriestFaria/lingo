package filters

import (
	"lingo/internal/analyzer/log"
)

type EmojiStrictFilter struct {}

func (f *EmojiStrictFilter) Apply(context *log.LogContext) []FilterIssue {
	return []FilterIssue{}
}