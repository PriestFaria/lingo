package analyzer

import (
	"go/ast"
	"lingo/internal/analyzer/log"
	"lingo/internal/config"
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

func analyzeMessage(pass *analysis.Pass, callExpr *ast.CallExpr, parts []log.LogPart, cfg *config.Config) {
	context := &log.LogContext{
		Pass:     pass,
		CallExpr: callExpr,
		Parts:    parts,
		FullText: buildFullText(parts),
	}

	var activeFilters []filters.LogFilter
	if cfg.Filters.IsEnabled("first_letter") {
		activeFilters = append(activeFilters, &filters.FirstLetterFilter{})
	}
	if cfg.Filters.IsEnabled("english") {
		activeFilters = append(activeFilters, &filters.EnglishFilter{})
	}
	if cfg.Filters.IsEnabled("emoji") {
		activeFilters = append(activeFilters, &filters.EmojiStrictFilter{})
	}
	if cfg.Filters.IsEnabled("security") {
		activeFilters = append(activeFilters, &filters.SecurityFilter{
			ExtraKeywords: cfg.Security.ExtraKeywords,
		})
	}

	pipeline := filters.NewFilterPipeline(activeFilters)

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
