package filters

import (
	"go/token"
	"unicode"
	"unicode/utf8"

	"github.com/PriestFaria/lingo/internal/analyzer/log"
)

// FirstLetterFilter reports log messages whose first literal part starts with
// an uppercase letter and provides an automatic fix to lowercase it.
type FirstLetterFilter struct{}

func (f *FirstLetterFilter) Apply(context *log.LogContext) []FilterIssue {
	for _, part := range context.Parts {
		if !part.IsLiteral || len(part.Value) == 0 {
			continue
		}

		firstRune, size := utf8.DecodeRuneInString(part.Value)
		if firstRune == utf8.RuneError {
			break
		}
		if !unicode.IsUpper(firstRune) {
			break 
		}

		contentStart := part.Pos + 1
		return []FilterIssue{{
			Message: "log message must start with a lowercase letter",
			Pos:     part.Pos,
			Fix: &IssueFix{
				Message: "lowercase first letter",
				Pos:     contentStart,
				End:     token.Pos(int(contentStart) + size),
				NewText: string(unicode.ToLower(firstRune)),
			},
		}}
	}
	return nil
}