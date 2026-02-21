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

func TestSecurityFilter_NoFalsePositive_Authenticated(t *testing.T) {
	f := &SecurityFilter{}
	ctx := makeCtx(makeParts("authenticated", false))
	issues := f.Apply(ctx)
	if len(issues) != 0 {
		t.Errorf("got %d issues, want 0 for 'authenticated' (should not match 'auth' as a substring)", len(issues))
	}
}

func TestSecurityFilter_SplitWords_LeadingUnderscore(t *testing.T) {
	f := &SecurityFilter{}
	ctx := makeCtx(makeParts("_password", false))
	issues := f.Apply(ctx)
	if len(issues) != 1 {
		t.Errorf("got %d issues, want 1 for '_password'", len(issues))
	}
}

func TestSecurityFilter_SplitWords_AllCapsAcronym(t *testing.T) {
	f := &SecurityFilter{}
	ctx := makeCtx(makeParts("APIKey", false))
	issues := f.Apply(ctx)
	if len(issues) != 1 {
		t.Errorf("got %d issues, want 1 for 'APIKey'", len(issues))
	}
}

func TestSecurityFilter_SplitWords_SnakeCase_PRIVATE_KEY(t *testing.T) {
	f := &SecurityFilter{}
	ctx := makeCtx(makeParts("PRIVATE_KEY", false))
	issues := f.Apply(ctx)
	if len(issues) != 1 {
		t.Errorf("got %d issues, want 1 for 'PRIVATE_KEY'", len(issues))
	}
}

func TestSecurityFilter_EmptyParts(t *testing.T) {
	f := &SecurityFilter{}
	ctx := makeCtx(makeParts("", true, "", false))
	issues := f.Apply(ctx)
	if len(issues) != 0 {
		t.Errorf("got %d issues, want 0 for empty parts", len(issues))
	}
}

func TestSecurityFilter_ExtraKeywords_LiteralMatch(t *testing.T) {
	f := &SecurityFilter{ExtraKeywords: []string{"cvv", "ssn"}}
	ctx := makeCtx(makeParts("processing cvv", true))
	issues := f.Apply(ctx)
	if len(issues) != 1 {
		t.Errorf("got %d issues, want 1 for extra keyword in literal", len(issues))
	}
}

func TestSecurityFilter_ExtraKeywords_VariableMatch(t *testing.T) {
	f := &SecurityFilter{ExtraKeywords: []string{"cvv", "ssn"}}
	ctx := makeCtx(makeParts("ssnNumber", false))
	issues := f.Apply(ctx)
	if len(issues) != 1 {
		t.Errorf("got %d issues, want 1 for extra keyword in variable name", len(issues))
	}
}

func TestSecurityFilter_ExtraKeywords_CaseInsensitive(t *testing.T) {
	f := &SecurityFilter{ExtraKeywords: []string{"OTP", "CVV"}}
	ctx := makeCtx(makeParts("otp code", true))
	issues := f.Apply(ctx)
	if len(issues) != 1 {
		t.Errorf("got %d issues, want 1 for case-insensitive extra keyword", len(issues))
	}
}

func TestSecurityFilter_ExtraKeywords_NoFalsePositive(t *testing.T) {
	f := &SecurityFilter{ExtraKeywords: []string{"ssn"}}
	ctx := makeCtx(makeParts("session", false))
	issues := f.Apply(ctx)
	if len(issues) != 0 {
		t.Errorf("got %d issues, want 0: 'session' should not match keyword 'ssn'", len(issues))
	}
}

func TestSecurityFilter_ExtraKeywords_MergedWithDefaults(t *testing.T) {
	f := &SecurityFilter{ExtraKeywords: []string{"cvv"}}
	ctx := makeCtx(makeParts("password", false))
	issues := f.Apply(ctx)
	if len(issues) != 1 {
		t.Errorf("got %d issues, want 1: built-in 'password' should still be detected", len(issues))
	}
}

func TestSecurityFilter_MultipleVariables(t *testing.T) {
	f := &SecurityFilter{}
	parts := makeParts(
		"user: ", true,
		"password", false,
		" key=", true,
		"apiKey", false,
	)
	ctx := makeCtx(parts)
	issues := f.Apply(ctx)
	if len(issues) != 3 {
		t.Errorf("got %d issues, want 3", len(issues))
	}
}
