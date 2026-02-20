package analyzer

import (
	"go/ast"
	"go/types"
	"lingo/internal/analyzer/log"
	"lingo/internal/filters"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func collectLogParts(expr ast.Expr, info *types.Info) []string {
	return []string{}
}
func analyzeMessage(pass *analysis.Pass, callExpr *ast.CallExpr, msgArg ast.Expr) {
	parts := collectLogParts(msgArg, pass.TypesInfo)
	fullText := strings.Join(parts, "")

	context := &log.LogContext{
		Pass: pass,
		CallExpr: callExpr,
		Parts: parts,
		FullText: fullText,
	}


	pipeline := filters.NewFilterPipeline([]filters.LogFilter{
		&filters.FirstLetterFilter{},
		&filters.EnglishFilter{},
		&filters.EmojiStrictFilter{},
	})

	pipeline.Process(context)
}