package filters

import (
	"fmt"
	"strings"
	"unicode"

	"lingo/internal/analyzer/log"
)

type SecurityFilter struct {
	ExtraKeywords []string
}

func (f *SecurityFilter) allKeywords() []string {
	if len(f.ExtraKeywords) == 0 {
		return sensitiveKeywords
	}
	all := make([]string, len(sensitiveKeywords), len(sensitiveKeywords)+len(f.ExtraKeywords))
	copy(all, sensitiveKeywords)
	for _, kw := range f.ExtraKeywords {
		all = append(all, strings.ToLower(kw))
	}
	return all
}

var sensitiveKeywords = []string{
	"password", "passwd", "pass",
	"secret",
	"token",
	"apikey", "api_key",
	"auth",
	"credential", "cred",
	"private", "privkey",
	"jwt",
	"key",
}

func splitWords(s string) []string {
	var words []string
	var cur strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if r == '_' || r == '-' {
			if cur.Len() > 0 {
				words = append(words, strings.ToLower(cur.String()))
				cur.Reset()
			}
			continue
		}
		if i > 0 && unicode.IsUpper(r) && unicode.IsLower(runes[i-1]) && cur.Len() > 0 {
			words = append(words, strings.ToLower(cur.String()))
			cur.Reset()
		}
		cur.WriteRune(r)
	}
	if cur.Len() > 0 {
		words = append(words, strings.ToLower(cur.String()))
	}
	return words
}

func containsSensitiveKeyword(name string, keywords []string) (string, bool) {
	words := splitWords(name)
	for _, word := range words {
		for _, kw := range keywords {
			if word == kw {
				return kw, true
			}
		}
	}
	return "", false
}

func containsSensitiveKeywordInLiteral(value string, keywords []string) (string, bool) {
	lower := strings.ToLower(value)
	words := strings.FieldsFunc(lower, func(r rune) bool {
		return r == ' ' || r == ':' || r == '=' || r == '_' || r == '-' || r == '/'
	})
	for _, word := range words {
		for _, kw := range keywords {
			if word == kw {
				return kw, true
			}
		}
	}
	return "", false
}

func (f *SecurityFilter) Apply(context *log.LogContext) []FilterIssue {
	keywords := f.allKeywords()
	var issues []FilterIssue
	for _, part := range context.Parts {
		if part.IsLiteral {
			if kw, ok := containsSensitiveKeywordInLiteral(part.Value, keywords); ok {
				issues = append(issues, FilterIssue{
					Message: fmt.Sprintf("log message may expose sensitive data: literal contains %q", kw),
					Pos:     part.Pos,
				})
			}
		} else {
			if kw, ok := containsSensitiveKeyword(part.Value, keywords); ok {
				issues = append(issues, FilterIssue{
					Message: fmt.Sprintf("log message may expose sensitive data: variable %q matches keyword %q", part.Value, kw),
					Pos:     part.Pos,
				})
			}
		}
	}
	return issues
}