package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"lingo/internal/config"
)

func TestLoad_EmptyPath_ReturnsDefault(t *testing.T) {
    cfg, err := config.Load("")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    for _, name := range []string{"first_letter", "english", "emoji", "security"} {
        if !cfg.Filters.IsEnabled(name) {
            t.Errorf("filter %q should be enabled by default", name)
        }
    }

    if len(cfg.Security.ExtraKeywords) != 0 {
        t.Errorf("expected no extra_keywords by default, got %v", cfg.Security.ExtraKeywords)
    }
}

func TestLoad_MissingFile_ReturnsError(t *testing.T) {
    _, err := config.Load("/nonexistent/path/.lingo.json")
    if err == nil {
        t.Fatal("expected error for missing file, got nil")
    }
}

func TestLoad_InvalidJSON_ReturnsError(t *testing.T) {
    f, err := os.CreateTemp(t.TempDir(), "*.lingo.json")
    if err != nil {
        t.Fatal(err)
    }
    f.WriteString("not { valid json")
    f.Close()

    _, err = config.Load(f.Name())
    if err == nil {
        t.Fatal("expected error for invalid JSON, got nil")
    }
}

func TestLoad_AllFiltersDisabled(t *testing.T) {
    content := `{
        "filters": {
            "first_letter": false,
            "english":      false,
            "emoji":        false,
            "security":     false
        }
    }`
    path := writeTemp(t, content)

    cfg, err := config.Load(path)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    for _, name := range []string{"first_letter", "english", "emoji", "security"} {
        if cfg.Filters.IsEnabled(name) {
            t.Errorf("filter %q should be disabled", name)
        }
    }
}

func TestLoad_PartialConfig_DefaultsApplied(t *testing.T) {
    content := `{"filters": {"first_letter": false}}`
    path := writeTemp(t, content)

    cfg, err := config.Load(path)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if cfg.Filters.IsEnabled("first_letter") {
        t.Error("first_letter should be disabled")
    }
    if !cfg.Filters.IsEnabled("english") {
        t.Error("english should be enabled (not specified → default)")
    }
    if !cfg.Filters.IsEnabled("security") {
        t.Error("security should be enabled (not specified → default)")
    }
}

func TestLoad_ExtraKeywords(t *testing.T) {
    content := `{
        "security": {
            "extra_keywords": ["cvv", "ssn", "OTP"]
        }
    }`
    path := writeTemp(t, content)

    cfg, err := config.Load(path)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if len(cfg.Security.ExtraKeywords) != 3 {
        t.Fatalf("expected 3 extra_keywords, got %d", len(cfg.Security.ExtraKeywords))
    }
}

func TestLoad_EmptyJSON_AllDefaults(t *testing.T) {
    path := writeTemp(t, `{}`)

    cfg, err := config.Load(path)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    for _, name := range []string{"first_letter", "english", "emoji", "security"} {
        if !cfg.Filters.IsEnabled(name) {
            t.Errorf("filter %q should be enabled when JSON is empty", name)
        }
    }
}

func TestFiltersConfig_IsEnabled_UnknownName(t *testing.T) {
    cfg := config.Default()
    if !cfg.Filters.IsEnabled("unknown_filter") {
        t.Error("unknown filter name should default to enabled")
    }
}

func TestDefault_RoundTrip(t *testing.T) {
    cfg := config.Default()
    data, err := json.Marshal(cfg)
    if err != nil {
        t.Fatalf("marshal failed: %v", err)
    }

    var restored config.Config
    if err := json.Unmarshal(data, &restored); err != nil {
        t.Fatalf("unmarshal failed: %v", err)
    }

    for _, name := range []string{"first_letter", "english", "emoji", "security"} {
        if !restored.Filters.IsEnabled(name) {
            t.Errorf("filter %q should be enabled after round-trip", name)
        }
    }
}

// writeTemp saves content to a temp file and returns its path.
func writeTemp(t *testing.T, content string) string {
    t.Helper()
    dir := t.TempDir()
    path := filepath.Join(dir, ".lingo.json")
    if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
        t.Fatal(err)
    }
    return path
}

// ── FromMap ─────────────────────────────────────────────────────────────────

func TestFromMap_Empty_ReturnsDefault(t *testing.T) {
    cfg, err := config.FromMap(map[string]any{})
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    for _, name := range []string{"first_letter", "english", "emoji", "security"} {
        if !cfg.Filters.IsEnabled(name) {
            t.Errorf("filter %q should be enabled by default", name)
        }
    }
    if len(cfg.Security.ExtraKeywords) != 0 {
        t.Errorf("expected no extra_keywords, got %v", cfg.Security.ExtraKeywords)
    }
}

func TestFromMap_InlineFilters(t *testing.T) {
    f := false
    _ = f
    cfg, err := config.FromMap(map[string]any{
        "filters": map[string]any{
            "first_letter": false,
            "english":      false,
        },
    })
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if cfg.Filters.IsEnabled("first_letter") {
        t.Error("first_letter should be disabled")
    }
    if cfg.Filters.IsEnabled("english") {
        t.Error("english should be disabled")
    }
    if !cfg.Filters.IsEnabled("emoji") {
        t.Error("emoji should be enabled (not specified → default)")
    }
    if !cfg.Filters.IsEnabled("security") {
        t.Error("security should be enabled (not specified → default)")
    }
}

func TestFromMap_InlineExtraKeywords(t *testing.T) {
    cfg, err := config.FromMap(map[string]any{
        "security": map[string]any{
            "extra_keywords": []any{"cvv", "ssn", "otp"},
        },
    })
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(cfg.Security.ExtraKeywords) != 3 {
        t.Fatalf("expected 3 extra_keywords, got %d", len(cfg.Security.ExtraKeywords))
    }
    for _, name := range []string{"first_letter", "english", "emoji", "security"} {
        if !cfg.Filters.IsEnabled(name) {
            t.Errorf("filter %q should be enabled by default", name)
        }
    }
}

func TestFromMap_AllFiltersDisabled(t *testing.T) {
    cfg, err := config.FromMap(map[string]any{
        "filters": map[string]any{
            "first_letter": false,
            "english":      false,
            "emoji":        false,
            "security":     false,
        },
    })
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    for _, name := range []string{"first_letter", "english", "emoji", "security"} {
        if cfg.Filters.IsEnabled(name) {
            t.Errorf("filter %q should be disabled", name)
        }
    }
}

func TestFromMap_UnknownKeysIgnored(t *testing.T) {
    cfg, err := config.FromMap(map[string]any{
        "unknown_key": "some_value",
        "another":     42,
    })
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    for _, name := range []string{"first_letter", "english", "emoji", "security"} {
        if !cfg.Filters.IsEnabled(name) {
            t.Errorf("filter %q should be enabled (unknown keys ignored)", name)
        }
    }
}