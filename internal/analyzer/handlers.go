package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"strings"

	"lingo/internal/analyzer/log"

	"golang.org/x/tools/go/analysis"
)

// collectPartsFromExpr рекурсивно обходит конкатенацию (+) и собирает LogPart.
// BasicLit → IsLiteral=true, Ident → IsLiteral=false.
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

func isFormatMethod(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), "f")
}

func collectArgs(callExpr *ast.CallExpr, info *types.Info) []log.LogPart {
	if len(callExpr.Args) == 0 {
		return nil
	}
	sel, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	var parts []log.LogPart
	if isFormatMethod(sel.Sel.Name) {
		// format string
		parts = append(parts, collectPartsFromExpr(callExpr.Args[0], info)...)
		// остальные аргументы — переменные
		for _, arg := range callExpr.Args[1:] {
			parts = append(parts, collectPartsFromExpr(arg, info)...)
		}
	} else {
		parts = collectPartsFromExpr(callExpr.Args[0], info)
	}
	return parts
}

func handleLog(pass *analysis.Pass, callExpr *ast.CallExpr) {
	parts := collectArgs(callExpr, pass.TypesInfo)
	if len(parts) == 0 {
		return
	}
	analyzeMessage(pass, callExpr, parts)
}

func handleSlog(pass *analysis.Pass, callExpr *ast.CallExpr) {
	if len(callExpr.Args) == 0 {
		return
	}
	parts := collectPartsFromExpr(callExpr.Args[0], pass.TypesInfo)
	if len(parts) == 0 {
		return
	}
	analyzeMessage(pass, callExpr, parts)
}

func handleZap(pass *analysis.Pass, callExpr *ast.CallExpr) {
	parts := collectArgs(callExpr, pass.TypesInfo)
	if len(parts) == 0 {
		return
	}
	analyzeMessage(pass, callExpr, parts)
}