package filters

import (
	"testing"
)

func TestEmojiStrictFilter_Emoji(t *testing.T) {
	f := &EmojiStrictFilter{}

	tests := []struct {
		name       string
		value      string
		isLiteral  bool
		wantIssues int
	}{
		{
			name:       "clean text — ok",
			value:      "server started",
			isLiteral:  true,
			wantIssues: 0,
		},
		{
			name:       "rocket emoji — issue",
			value:      "server started 🚀",
			isLiteral:  true,
			wantIssues: 1,
		},
		{
			name:       "fire emoji — issue",
			value:      "🔥 error occurred",
			isLiteral:  true,
			wantIssues: 1,
		},
		{
			name:       "double exclamation — issue",
			value:      "connection failed!!",
			isLiteral:  true,
			wantIssues: 1,
		},
		{
			name:       "triple question — issue",
			value:      "what???",
			isLiteral:  true,
			wantIssues: 1,
		},
		{
			name:       "ellipsis — issue",
			value:      "loading...",
			isLiteral:  true,
			wantIssues: 1,
		},
		{
			name:       "single exclamation — ok",
			value:      "connection failed!",
			isLiteral:  true,
			wantIssues: 0,
		},
		{
			name:       "single dot — ok",
			value:      "something went wrong.",
			isLiteral:  true,
			wantIssues: 0,
		},
		{
			name:       "non-literal with emoji — ok (variables are not checked)",
			value:      "🚀",
			isLiteral:  false,
			wantIssues: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := makeCtx(makeParts(tc.value, tc.isLiteral))
			issues := f.Apply(ctx)
			if len(issues) != tc.wantIssues {
				t.Errorf("got %d issues, want %d (value=%q)", len(issues), tc.wantIssues, tc.value)
			}
		})
	}
}

func TestEmojiStrictFilter_EmojiAndRepeatedPunct(t *testing.T) {
	f := &EmojiStrictFilter{}
	ctx := makeCtx(makeParts("error!!! 🚀", true))
	issues := f.Apply(ctx)
	if len(issues) != 2 {
		t.Errorf("got %d issues, want 2", len(issues))
	}
}
