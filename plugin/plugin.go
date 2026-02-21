// Package main is the golangci-lint plugin entry point.
//
// Build the plugin:
//
//	go build -buildmode=plugin -o lingo.so ./plugin/
//
// Configure golangci-lint (.golangci.yml):
//
//	linters:
//	  enable:
//	    - lingo
//	linters-settings:
//	  custom:
//	    lingo:
//	      type: goplugin
//	      path: ./lingo.so
//	      description: "Linter for log messages: style, language and sensitive data"
//	      settings:
//	        config: .lingo.json   # optional path to .lingo.json
//
// The "config" key in settings is optional. When omitted, all filters are
// enabled with default settings (no extra security keywords).
package main

import (
	"fmt"

	"golang.org/x/tools/go/analysis"

	"lingo/internal/analyzer"
	"lingo/internal/config"
)

// New is the constructor called by golangci-lint's Go plugin system.
//
// conf may be a map[string]any. Supported keys:
//   - "config" (string): path to a .lingo.json config file.
//
// When conf is nil or "config" is absent, the default config is used
// (all filters enabled, no extra security keywords).
func New(conf any) ([]*analysis.Analyzer, error) {
	cfg := config.Default()

	if m, ok := conf.(map[string]any); ok {
		if path, ok := m["config"].(string); ok && path != "" {
			loaded, err := config.Load(path)
			if err != nil {
				return nil, fmt.Errorf("lingo: load config %q: %w", path, err)
			}
			cfg = loaded
		}
	}

	return []*analysis.Analyzer{analyzer.NewAnalyzerWithConfig(cfg)}, nil
}
