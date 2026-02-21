package analyzer_test

import (
	"path/filepath"
	"testing"

	"github.com/PriestFaria/lingo/internal/analyzer"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.Analyzer, "basic", "withzap", "clean", "concat", "realworld")
}

func TestAnalyzerWithConfig(t *testing.T) {
	testdata := analysistest.TestData()
	configFile := filepath.Join(testdata, "src", "withconfig", ".lingo.json")

	if err := analyzer.Analyzer.Flags.Set("config", configFile); err != nil {
		t.Fatalf("failed to set config flag: %v", err)
	}
	t.Cleanup(func() {
		analyzer.Analyzer.Flags.Set("config", "") //nolint:errcheck
	})

	analysistest.Run(t, testdata, analyzer.Analyzer, "withconfig")
}
