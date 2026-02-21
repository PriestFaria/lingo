package analyzer

import (
	"go/ast"

	"lingo/internal/config"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var configPath string

func init() {
	Analyzer.Flags.StringVar(&configPath, "config", "", "path to lingo config file (.lingo.json)")
}

var Analyzer *analysis.Analyzer = &analysis.Analyzer{
	Name: "lingo",
	Doc:  "Lingo is a static analysis tool that detects logs issues, such as first-letter-case, language strict and leaking business information.",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}

	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		callExpession, ok := n.(*ast.CallExpr)
		if !ok {
			return
		}

		selectorExpression, ok := callExpession.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		selection, ok := pass.TypesInfo.Selections[selectorExpression]
		var pkgPath string

		if ok {
			pkgPath = selection.Obj().Pkg().Path()
		} else if obj, ok := pass.TypesInfo.Uses[selectorExpression.Sel]; ok {
			if pkg := obj.Pkg(); pkg != nil {
				pkgPath = pkg.Path()
			}
		}

		switch pkgPath {
		case "go.uber.org/zap":
			handleZap(pass, callExpession, cfg)
		case "log/slog":
			handleSlog(pass, callExpession, cfg)
		case "log":
			handleLog(pass, callExpession, cfg)
		}
	})
	return nil, nil
}
