package analyzer

import (
	"go/ast"
	"github.com/PriestFaria/lingo/internal/analyzer/log"
	"github.com/PriestFaria/lingo/internal/config"
	"github.com/PriestFaria/lingo/internal/filters"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// buildFullText concatenates the Value fields of all parts into a single string.
func buildFullText(parts []log.LogPart) string {
	var sb strings.Builder
	for _, p := range parts {
		sb.WriteString(p.Value)
	}
	return sb.String()
}

// analyzeMessage builds a LogContext from parts, constructs the filter pipeline
// according to cfg, and reports any issues found via pass.Report/pass.Reportf.
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
