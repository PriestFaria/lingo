package filters

import (
	"testing"
)

func TestFirstLetterFilter(t *testing.T) {
	f := &FirstLetterFilter{}

	tests := []struct {
		name        string
		parts       []string 
		wantIssues  int
		wantFix     bool
	}{
		{
			name:       "lowercase — ok",
			parts:      []string{"starting server"},
			wantIssues: 0,
		},
		{
			name:       "uppercase — issue with fix",
			parts:      []string{"Starting server"},
			wantIssues: 1,
			wantFix:    true,
		},
		{
			name:       "empty literal — ok",
			parts:      []string{""},
			wantIssues: 0,
		},
		{
			name:       "first part lowercase — ok (non-literal second)",
			parts:      []string{"connecting to"},
			wantIssues: 0,
		},
		{
			name:       "single uppercase letter",
			parts:      []string{"E"},
			wantIssues: 1,
			wantFix:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var literalParts []interface{}
			for _, v := range tc.parts {
				literalParts = append(literalParts, v, true)
			}
			ctx := makeCtx(makeParts(literalParts...))
			issues := f.Apply(ctx)

			if len(issues) != tc.wantIssues {
				t.Errorf("got %d issues, want %d", len(issues), tc.wantIssues)
			}
			if tc.wantIssues > 0 && tc.wantFix && issues[0].Fix == nil {
				t.Errorf("expected a Fix, got nil")
			}
			if tc.wantIssues > 0 && !tc.wantFix && issues[0].Fix != nil {
				t.Errorf("expected no Fix, got one")
			}
		})
	}
}

func TestFirstLetterFilter_SkipsNonLiteral(t *testing.T) {
	f := &FirstLetterFilter{}

	parts := makeParts("someVar", false, "Starting server", true)
	ctx := makeCtx(parts)
	issues := f.Apply(ctx)
	if len(issues) != 1 {
		t.Errorf("got %d issues, want 1", len(issues))
	}
}

func TestFirstLetterFilter_FixValue(t *testing.T) {
	f := &FirstLetterFilter{}
	ctx := makeCtx(makeParts("Failed to connect", true))
	issues := f.Apply(ctx)
	if len(issues) != 1 {
		t.Fatalf("got %d issues, want 1", len(issues))
	}
	if issues[0].Fix == nil {
		t.Fatal("expected Fix, got nil")
	}
	if issues[0].Fix.NewText != "f" {
		t.Errorf("Fix.NewText = %q, want %q", issues[0].Fix.NewText, "f")
	}
}
