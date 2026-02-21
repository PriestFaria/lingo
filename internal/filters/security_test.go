package filters

import (
	"testing"
)

func TestSecurityFilter(t *testing.T) {
	f := &SecurityFilter{}

	tests := []struct {
		name       string
		value      string
		isLiteral  bool
		wantIssues int
	}{
		{
			name:       "clean literal — ok",
			value:      "user authenticated successfully",
			isLiteral:  true,
			wantIssues: 0,
		},
		{
			name:       "clean variable — ok",
			value:      "userID",
			isLiteral:  false,
			wantIssues: 0,
		},
		{
			name:       "variable named password — issue",
			value:      "password",
			isLiteral:  false,
			wantIssues: 1,
		},
		{
			name:       "variable named token — issue",
			value:      "token",
			isLiteral:  false,
			wantIssues: 1,
		},
		{
			name:       "variable named apiKey — issue",
			value:      "apiKey",
			isLiteral:  false,
			wantIssues: 1,
		},
		{
			name:       "variable named userPassword — issue (содержит keyword)",
			value:      "userPassword",
			isLiteral:  false,
			wantIssues: 1,
		},
		{
			name:       "variable named jwtToken — issue",
			value:      "jwtToken",
			isLiteral:  false,
			wantIssues: 1,
		},
		{
			name:       "literal with password: marker — issue",
			value:      "user password: ",
			isLiteral:  true,
			wantIssues: 1,
		},
		{
			name:       "literal with token= marker — issue",
			value:      "token=",
			isLiteral:  true,
			wantIssues: 1,
		},
		{
			name:       "literal with api_key — issue",
			value:      "api_key=",
			isLiteral:  true,
			wantIssues: 1,
		},
		{
			name:       "variable named PASSWORD (uppercase) — issue",
			value:      "PASSWORD",
			isLiteral:  false,
			wantIssues: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := makeCtx(makeParts(tc.value, tc.isLiteral))
			issues := f.Apply(ctx)
			if len(issues) != tc.wantIssues {
				t.Errorf("got %d issues, want %d (value=%q, isLiteral=%v)",
					len(issues), tc.wantIssues, tc.value, tc.isLiteral)
			}
		})
	}
}

func TestSecurityFilter_MultipleVariables(t *testing.T) {
	f := &SecurityFilter{}
	// log.Info("user: " + password + " key=" + apiKey)
	parts := makeParts(
		"user: ", true,
		"password", false,
		" key=", true,
		"apiKey", false,
	)
	ctx := makeCtx(parts)
	issues := f.Apply(ctx)
	// password (var) + " key=" (literal marker) + apiKey (var) = 3
	if len(issues) != 3 {
		t.Errorf("got %d issues, want 3", len(issues))
	}
}
