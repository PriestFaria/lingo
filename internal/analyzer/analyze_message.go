package analyzer

import (
	"go/ast"
	"lingo/internal/analyzer/log"
	"lingo/internal/filters"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func buildFullText(parts []log.LogPart) string {
	var sb strings.Builder
	for _, p := range parts {
		sb.WriteString(p.Value)
	}
	return sb.String()
}

func analyzeMessage(pass *analysis.Pass, callExpr *ast.CallExpr, parts []log.LogPart) {
	context := &log.LogContext{
		Pass:     pass,
		CallExpr: callExpr,
		Parts:    parts,
		FullText: buildFullText(parts),
	}

	pipeline := filters.NewFilterPipeline([]filters.LogFilter{
		&filters.FirstLetterFilter{},
		&filters.EnglishFilter{},
		&filters.EmojiStrictFilter{},
		&filters.SecurityFilter{},
	})

	issues := pipeline.Process(context)
	for _, issue := range issues {
		if issue.Fix != nil {
			pass.Report(analysis.Diagnostic{
				Pos:     issue.Pos,
				Message: issue.Message,
				SuggestedFixes: []analysis.SuggestedFix{{
					Message: issue.Fix.Message,
					TextEdits: []analysis.TextEdit{{
						Pos:     issue.Fix.Pos,
						End:     issue.Fix.End,
						NewText: []byte(issue.Fix.NewText),
					}},
				}},
			})
		} else {
			pass.Reportf(issue.Pos, "%s", issue.Message)
		}
	}
}
