package filters

import "go/token"

// IssueFix describes a single-span text edit that resolves a FilterIssue.
// It is applied by the analyzer as an analysis.SuggestedFix.
type IssueFix struct {
	// Message is the human-readable description shown in the fix suggestion.
	Message string
	// Pos is the start position of the text to be replaced (inclusive).
	Pos token.Pos
	// End is the end position of the text to be replaced (exclusive).
	End token.Pos
	// NewText is the replacement string.
	NewText string
}

// FilterIssue represents a single violation found by a LogFilter.
// Fix is nil when no automatic fix is available.
type FilterIssue struct {
	Message string
	Pos     token.Pos
	Fix     *IssueFix
}