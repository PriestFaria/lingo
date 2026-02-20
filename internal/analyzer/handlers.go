package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)


func handleZap(pass *analysis.Pass, callExpr *ast.CallExpr) {}


func handleSlog(pass *analysis.Pass, callExpr *ast.CallExpr) {}


func handleLog(pass *analysis.Pass, callExpr *ast.CallExpr) {}