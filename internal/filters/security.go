package filters

import (
	"fmt"
	"strings"

	"lingo/internal/analyzer/log"
)

type SecurityFilter struct{}

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

func containsSensitiveKeyword(name string) (string, bool) {
	lower := strings.ToLower(name)
	for _, kw := range sensitiveKeywords {
		if strings.Contains(lower, kw) {
			return kw, true
		}
	}
	return "", false
}

func (f *SecurityFilter) Apply(context *log.LogContext) []FilterIssue {
	var issues []FilterIssue
	for _, part := range context.Parts {
		if part.IsLiteral {
			if kw, ok := containsSensitiveKeyword(part.Value); ok {
				issues = append(issues, FilterIssue{
					Message: fmt.Sprintf("log message may expose sensitive data: literal contains %q", kw),
					Pos:     part.Pos,
				})
			}
		} else {
			if kw, ok := containsSensitiveKeyword(part.Value); ok {
				issues = append(issues, FilterIssue{
					Message: fmt.Sprintf("log message may expose sensitive data: variable %q matches keyword %q", part.Value, kw),
					Pos:     part.Pos,
				})
			}
		}
	}
	return issues
}