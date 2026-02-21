// Package main is the golangci-lint plugin entry point.
//
// Build the plugin:
//
//	go build -buildmode=plugin -o lingo.so ./plugin/
//
// Configure golangci-lint (.golangci.yml) — inline config (recommended):
//
//	linters:
//	  enable:
//	    - lingo
//	linters-settings:
//	  custom:
//	    lingo:
//	      type: goplugin
//	      path: ./lingo.so
//	      settings:
//	        filters:
//	          first_letter: true
//	          english: true
//	          emoji: true
//	          security: true
//	        security:
//	          extra_keywords:
//	            - cvv
//	            - ssn
//	            - otp
//
// Alternatively, point to an external .lingo.json file:
//
//	settings:
//	  config: .lingo.json
//
// Priority: inline (filters/security keys) > config file > defaults.
// When settings is omitted entirely, all filters are enabled with no extra keywords.
package main

import (
	"fmt"

	"golang.org/x/tools/go/analysis"

	"github.com/PriestFaria/lingo/internal/analyzer"
	"github.com/PriestFaria/lingo/internal/config"
)

// New is the constructor called by golangci-lint's Go plugin system.
//
// conf may be a map[string]any built from the golangci-lint settings block.
//
// Resolution priority:
//  1. Inline — if "filters" or "security" keys are present, the map is parsed
//     directly into Config (same structure as .lingo.json).
//  2. File   — if "config" key (string) is present, the file is loaded.
//  3. Default — all filters enabled, no extra keywords.
func New(conf any) ([]*analysis.Analyzer, error) {
	cfg, err := resolveConfig(conf)
	if err != nil {
		return nil, err
	}
	return []*analysis.Analyzer{analyzer.NewAnalyzerWithConfig(cfg)}, nil
}

func resolveConfig(conf any) (*config.Config, error) {
	m, ok := conf.(map[string]any)
	if !ok || len(m) == 0 {
		return config.Default(), nil
	}

	// Inline config: filters and/or security keys present directly in settings.
	_, hasFilters := m["filters"]
	_, hasSecurity := m["security"]
	if hasFilters || hasSecurity {
		cfg, err := config.FromMap(m)
		if err != nil {
			return nil, fmt.Errorf("lingo: parse inline settings: %w", err)
		}
		return cfg, nil
	}

	// File fallback: "config" key points to a .lingo.json path.
	if path, ok := m["config"].(string); ok && path != "" {
		cfg, err := config.Load(path)
		if err != nil {
			return nil, fmt.Errorf("lingo: load config %q: %w", path, err)
		}
		return cfg, nil
	}

	return config.Default(), nil
}
