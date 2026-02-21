package filters

import (
	"testing"
)

func TestEmojiStrictFilter_Emoji(t *testing.T) {
	f := &EmojiStrictFilter{}

	tests := []struct {
		name       string
		value      string
		wantIssues int
	}{
		{
			name:       "clean text — ok",
			value:      "server started",
			wantIssues: 0,
		},
		{
			name:       "rocket emoji — issue",
			value:      "server started 🚀",
			wantIssues: 1,
		},
		{
			name:       "fire emoji — issue",
			value:      "🔥 error occurred",
			wantIssues: 1,
		},
		{
			name:       "double exclamation — issue",
			value:      "connection failed!!",
			wantIssues: 1,
		},
		{
			name:       "triple question — issue",
			value:      "what???",
			wantIssues: 1,
		},
		{
			name:       "ellipsis — issue",
			value:      "loading...",
			wantIssues: 1,
		},
		{
			name:       "single exclamation — ok",
			value:      "connection failed!",
			wantIssues: 0,
		},
		{
			name:       "single dot — ok",
			value:      "something went wrong.",
			wantIssues: 0,
		},
		{
			name:       "non-literal with emoji — ok (не проверяем переменные)",
			value:      "🚀",
			wantIssues: 0, 
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			isLiteral := tc.name != "non-literal with emoji — ok (не проверяем переменные)"
			ctx := makeCtx(makeParts(tc.value, isLiteral))
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
