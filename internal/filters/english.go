package filters

import (
	"fmt"
	"unicode"

	"lingo/internal/analyzer/log"
)

// EnglishFilter reports log messages that contain non-ASCII letters,
// indicating the message is not written in English.
type EnglishFilter struct{}

func (f *EnglishFilter) Apply(context *log.LogContext) []FilterIssue {
	var issues []FilterIssue
	for _, part := range context.Parts {
		if !part.IsLiteral {
			continue
		}
		for _, r := range part.Value {
			if r > 127 && unicode.IsLetter(r) {
				issues = append(issues, FilterIssue{
					Message: fmt.Sprintf("log message must be in English, found non-ASCII character: %q", r),
					Pos:     part.Pos,
					Fix:     nil,
				})
				break 
			}
		}
	}
	return issues
}