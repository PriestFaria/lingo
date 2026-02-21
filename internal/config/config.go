package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// FiltersConfig manages enabling/disabling of individual filters.
// A nil *bool means "not configured" and defaults to enabled.
type FiltersConfig struct {
    FirstLetter *bool `json:"first_letter"`
    English     *bool `json:"english"`
    Emoji       *bool `json:"emoji"`
    Security    *bool `json:"security"`
}

// IsEnabled returns true if the named filter is enabled.
// Recognised names: "first_letter", "english", "emoji", "security".
func (f *FiltersConfig) IsEnabled(name string) bool {
    var p *bool
    switch name {
    case "first_letter":
        p = f.FirstLetter
    case "english":
        p = f.English
    case "emoji":
        p = f.Emoji
    case "security":
        p = f.Security
    }
    return p == nil || *p
}

// SecurityConfig holds settings for SecurityFilter.
type SecurityConfig struct {
    // ExtraKeywords are additional sensitive keywords beyond the built-in list.
    // Matching is case-insensitive: "CVV" and "cvv" are equivalent.
    ExtraKeywords []string `json:"extra_keywords"`
}

// Config is the root configuration structure for a .lingo.json file.
//
// Example .lingo.json:
//
//	{
//	  "filters": { "first_letter": false },
//	  "security": { "extra_keywords": ["cvv", "ssn"] }
//	}
type Config struct {
    Filters  FiltersConfig  `json:"filters"`
    Security SecurityConfig `json:"security"`
}

// Default returns the default configuration:
// all filters enabled, no custom keywords.
func Default() *Config {
    return &Config{}
}

// FromMap builds a Config from a map[string]any (e.g. from golangci-lint settings).
// The map structure mirrors the JSON layout of .lingo.json:
//
//	map[string]any{
//	    "filters": map[string]any{
//	        "first_letter": false,
//	        "english":      true,
//	    },
//	    "security": map[string]any{
//	        "extra_keywords": []any{"cvv", "ssn"},
//	    },
//	}
//
// Unknown keys are silently ignored. An empty map returns Default().
func FromMap(m map[string]any) (*Config, error) {
	if len(m) == 0 {
		return Default(), nil
	}

	data, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("lingo: cannot encode settings map: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("lingo: cannot parse settings map: %w", err)
	}

	return &cfg, nil
}

// Load reads a Config from a JSON file at path.
// If path is empty, Default() is returned.
// Fields absent from the file keep their zero value (filter enabled).
func Load(path string) (*Config, error) {
    if path == "" {
        return Default(), nil
    }

    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("lingo: cannot read config %q: %w", path, err)
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("lingo: cannot parse config %q: %w", path, err)
    }

    return &cfg, nil
}