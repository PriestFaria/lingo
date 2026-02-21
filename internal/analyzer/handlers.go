package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"strings"

	"github.com/PriestFaria/lingo/internal/analyzer/log"
	"github.com/PriestFaria/lingo/internal/config"

	"golang.org/x/tools/go/analysis"
)

// collectPartsFromExpr recursively decomposes an AST expression into LogParts.
// It handles string concatenation (BinaryExpr with ADD), string literals
// (BasicLit), and identifiers (Ident).
func collectPartsFromExpr(expr ast.Expr, info *types.Info) []log.LogPart {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			left := collectPartsFromExpr(e.X, info)
			right := collectPartsFromExpr(e.Y, info)
			return append(left, right...)
		}
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			value, err := strconv.Unquote(e.Value)
			if err != nil {
				value = strings.Trim(e.Value, "`")
			}
			return []log.LogPart{{
				Value:     value,
				IsLiteral: true,
				Pos:       e.Pos(),
				End:       e.End(),
			}}
		}
	case *ast.Ident:
		return []log.LogPart{{
			Value:     e.Name,
			IsLiteral: false,
			Pos:       e.Pos(),
			End:       e.End(),
		}}
	}
	return nil
}

// isFormatMethod reports whether name is a format-style method (e.g. Printf,
// Infof) by checking for a trailing "f" suffix (case-insensitive).
func isFormatMethod(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), "f")
}

// collectArgs extracts LogParts from a call expression.
// For format methods the format string (Args[0]) and all value arguments are
// collected; for regular methods only the first argument is used.
func collectArgs(callExpr *ast.CallExpr, info *types.Info) []log.LogPart {
	if len(callExpr.Args) == 0 {
		return nil
	}
	sel := callExpr.Fun.(*ast.SelectorExpr)

	var parts []log.LogPart
	if isFormatMethod(sel.Sel.Name) {
		parts = append(parts, collectPartsFromExpr(callExpr.Args[0], info)...)
		for _, arg := range callExpr.Args[1:] {
			parts = append(parts, collectPartsFromExpr(arg, info)...)
		}
	} else {
		parts = collectPartsFromExpr(callExpr.Args[0], info)
	}
	return parts
}

// handleLog processes a call to the standard library "log" package.
func handleLog(pass *analysis.Pass, callExpr *ast.CallExpr, cfg *config.Config) {
	parts := collectArgs(callExpr, pass.TypesInfo)
	if len(parts) == 0 {
		return
	}
	analyzeMessage(pass, callExpr, parts, cfg)
}

// handleSlog processes a call to the "log/slog" package.
// Only the first argument (the message string) is inspected.
func handleSlog(pass *analysis.Pass, callExpr *ast.CallExpr, cfg *config.Config) {
	if len(callExpr.Args) == 0 {
		return
	}
	parts := collectPartsFromExpr(callExpr.Args[0], pass.TypesInfo)
	if len(parts) == 0 {
		return
	}
	analyzeMessage(pass, callExpr, parts, cfg)
}

// handleZap processes a call to "go.uber.org/zap".
func handleZap(pass *analysis.Pass, callExpr *ast.CallExpr, cfg *config.Config) {
	parts := collectArgs(callExpr, pass.TypesInfo)
	if len(parts) == 0 {
		return
	}
	analyzeMessage(pass, callExpr, parts, cfg)
}