package filters

import (
	"testing"
)

func TestEnglishFilter(t *testing.T) {
	f := &EnglishFilter{}

	tests := []struct {
		name       string
		value      string
		isLiteral  bool
		wantIssues int
	}{
		{
			name:       "english text — ok",
			value:      "starting server on port 8080",
			isLiteral:  true,
			wantIssues: 0,
		},
		{
			name:       "cyrillic — issue",
			value:      "запуск сервера",
			isLiteral:  true,
			wantIssues: 1,
		},
		{
			name:       "mixed english and cyrillic — issue",
			value:      "server запущен",
			isLiteral:  true,
			wantIssues: 1,
		},
		{
			name:       "non-literal cyrillic — ok (variables are not checked)",
			value:      "кириллица",
			isLiteral:  false,
			wantIssues: 0,
		},
		{
			name:       "digits and punctuation — ok",
			value:      "connected on port 8080!",
			isLiteral:  true,
			wantIssues: 0,
		},
		{
			name:       "empty string — ok",
			value:      "",
			isLiteral:  true,
			wantIssues: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := makeCtx(makeParts(tc.value, tc.isLiteral))
			issues := f.Apply(ctx)
			if len(issues) != tc.wantIssues {
				t.Errorf("got %d issues, want %d", len(issues), tc.wantIssues)
			}
		})
	}
}

// TestEnglishFilter_OncePerPart verifies that at most one issue is reported per
// part, regardless of how many non-ASCII characters are present.
func TestEnglishFilter_OncePerPart(t *testing.T) {
	f := &EnglishFilter{}
	ctx := makeCtx(makeParts("привет мир", true))
	issues := f.Apply(ctx)
	if len(issues) != 1 {
		t.Errorf("got %d issues, want exactly 1 per part", len(issues))
	}
}

// TestEnglishFilter_MultiplePartsWithIssues verifies that each violating part
// produces exactly one issue.
func TestEnglishFilter_MultiplePartsWithIssues(t *testing.T) {
	f := &EnglishFilter{}
	ctx := makeCtx(makeParts("привет", true, "мир", true))
	issues := f.Apply(ctx)
	if len(issues) != 2 {
		t.Errorf("got %d issues, want 2", len(issues))
	}
}
