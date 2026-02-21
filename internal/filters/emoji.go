package filters

import (
	"fmt"
	"regexp"

	"lingo/internal/analyzer/log"
)

// EmojiStrictFilter reports log messages that contain emoji characters or
// repeated punctuation sequences (e.g. !!, ???, ...).
type EmojiStrictFilter struct{}

// repeatedPunct matches two or more consecutive ! or ? and two or more dots.
var repeatedPunct = regexp.MustCompile(`[!?]{2,}|\.{2,}`)

// isEmoji reports whether the rune falls within a known emoji Unicode range.
func isEmoji(r rune) bool {
	return (r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
		(r >= 0x1F300 && r <= 0x1F5FF) || // Misc Symbols & Pictographs
		(r >= 0x1F680 && r <= 0x1F6FF) || // Transport & Map
		(r >= 0x1F1E0 && r <= 0x1F1FF) || // Flags
		(r >= 0x2600 && r <= 0x26FF) || // Misc Symbols
		(r >= 0x2700 && r <= 0x27BF) // Dingbats
}

func (f *EmojiStrictFilter) Apply(context *log.LogContext) []FilterIssue {
	var issues []FilterIssue
	for _, part := range context.Parts {
		if !part.IsLiteral {
			continue
		}

		for _, r := range part.Value {
			if isEmoji(r) {
				issues = append(issues, FilterIssue{
					Message: fmt.Sprintf("log message must not contain emoji: %q", r),
					Pos:     part.Pos,
				})
				break
			}
		}

		if loc := repeatedPunct.FindStringIndex(part.Value); loc != nil {
			issues = append(issues, FilterIssue{
				Message: fmt.Sprintf("log message must not contain repeated punctuation: %q", part.Value[loc[0]:loc[1]]),
				Pos:     part.Pos,
			})
		}
	}
	return issues
}