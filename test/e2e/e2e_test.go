//go:build e2e

// Package e2e tests the full linter lifecycle:
// build cmd/lingo binary → run via go vet → verify diagnostics.
//
// Run with: go test -tags e2e ./test/e2e/
package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func projectRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// test/e2e/ → two levels up to the repository root
	return filepath.Join(wd, "..", "..")
}

func buildLingo(t *testing.T) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "lingo-check")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	root := projectRoot(t)
	cmd := exec.Command("go", "build", "-o", binary, filepath.Join(root, "cmd", "lingo"))
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build lingo: %s\n%s", err, out)
	}
	return binary
}

// runVet runs go vet without a config flag in the given directory.
func runVet(t *testing.T, binary, projectDir string) (string, error) {
	t.Helper()
	cmd := exec.Command("go", "vet", "-vettool="+binary, "./...")
	cmd.Dir = projectDir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// runVetWithConfig runs go vet with the -config flag.
// configPath must be an absolute path to a .lingo.json file.
func runVetWithConfig(t *testing.T, binary, projectDir, configPath string) (string, error) {
	t.Helper()
	cmd := exec.Command("go", "vet", "-vettool="+binary, "-config="+configPath, "./...")
	cmd.Dir = projectDir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// absConfig returns the absolute path to the .lingo.json file inside a testdata project.
func absConfig(t *testing.T, projectName string) string {
	t.Helper()
	path, err := filepath.Abs(filepath.Join("testdata", projectName, ".lingo.json"))
	if err != nil {
		t.Fatal(err)
	}
	return path
}

// ── Base tests ───────────────────────────────────────────────────────────────

// TestE2E_CleanProject verifies that the linter reports no issues on clean code.
func TestE2E_CleanProject(t *testing.T) {
	binary := buildLingo(t)
	projectDir := filepath.Join("testdata", "clean-project")

	out, err := runVet(t, binary, projectDir)
	if err != nil {
		t.Errorf("lingo reported issues on clean code (expected no issues):\n%s", out)
	}
}

// TestE2E_ViolationsProject verifies that the linter detects all violation types.
func TestE2E_ViolationsProject(t *testing.T) {
	binary := buildLingo(t)
	projectDir := filepath.Join("testdata", "violations-project")

	out, err := runVet(t, binary, projectDir)
	if err == nil {
		t.Fatal("lingo found no issues in violations project (expected several)")
	}

	expected := []string{
		"must start with a lowercase letter",
		"must be in English",
		"must not contain emoji",
		"must not contain repeated punctuation",
		"may expose sensitive data",
	}
	for _, msg := range expected {
		if !strings.Contains(out, msg) {
			t.Errorf("expected diagnostic not found: %q\nfull output:\n%s", msg, out)
		}
	}
}

// ── Config tests ─────────────────────────────────────────────────────────────

// TestE2E_Config_DisabledFilters verifies that when all filters are disabled
// in the config the linter reports no issues, even on code with violations.
func TestE2E_Config_DisabledFilters(t *testing.T) {
	binary := buildLingo(t)
	projectDir := filepath.Join("testdata", "disabled-filters-project")
	configPath := absConfig(t, "disabled-filters-project")

	out, err := runVetWithConfig(t, binary, projectDir, configPath)
	if err != nil {
		t.Errorf("expected no issues when all filters disabled, but got:\n%s", out)
	}
}

// TestE2E_Config_DisabledFilters_WithoutConfig runs the same violation code
// without a config and expects errors, because all filters are enabled by default.
func TestE2E_Config_DisabledFilters_WithoutConfig(t *testing.T) {
	binary := buildLingo(t)
	projectDir := filepath.Join("testdata", "disabled-filters-project")

	out, err := runVet(t, binary, projectDir)
	if err == nil {
		t.Fatalf("expected issues without config, but lingo reported nothing\noutput:\n%s", out)
	}
}

// TestE2E_Config_ExtraKeywords verifies that custom keywords (cvv, ssn, otp)
// configured via extra_keywords are detected alongside built-in ones (password).
func TestE2E_Config_ExtraKeywords(t *testing.T) {
	binary := buildLingo(t)
	projectDir := filepath.Join("testdata", "extra-keywords-project")
	configPath := absConfig(t, "extra-keywords-project")

	out, err := runVetWithConfig(t, binary, projectDir, configPath)
	if err == nil {
		t.Fatalf("expected security issues but lingo found nothing\noutput:\n%s", out)
	}

	// every custom keyword must be detected
	expectedKeywords := []string{"cvv", "ssn", "otp"}
	for _, kw := range expectedKeywords {
		if !strings.Contains(out, kw) {
			t.Errorf("expected report for custom keyword %q not found\nfull output:\n%s", kw, out)
		}
	}

	// built-in keyword must still be detected
	if !strings.Contains(out, "password") {
		t.Errorf("expected report for built-in keyword %q not found\nfull output:\n%s", "password", out)
	}

	// safe variable requestID must not trigger a false positive
	if strings.Contains(out, "requestID") {
		t.Errorf("false positive: requestID should not trigger security filter\nfull output:\n%s", out)
	}

	// all security diagnostics must carry the standard message
	if !strings.Contains(out, "may expose sensitive data") {
		t.Errorf("expected standard security message not found\nfull output:\n%s", out)
	}
}

// TestE2E_Config_ExtraKeywords_NotDetectedWithoutConfig verifies that without
// a config the custom keywords (cvv, ssn, otp) are unknown to the linter and
// do not trigger any diagnostics.
func TestE2E_Config_ExtraKeywords_NotDetectedWithoutConfig(t *testing.T) {
	binary := buildLingo(t)
	projectDir := filepath.Join("testdata", "extra-keywords-project")

	out, err := runVet(t, binary, projectDir)

	// built-in "password" must still be detected
	if err == nil {
		t.Fatalf("expected at least one issue (password), but lingo found nothing\noutput:\n%s", out)
	}
	if !strings.Contains(out, "password") {
		t.Errorf("expected built-in 'password' detection\nfull output:\n%s", out)
	}

	// custom keywords must NOT fire without a config
	for _, kw := range []string{"\"cvv\"", "\"ssn\"", "\"otp\""} {
		if strings.Contains(out, kw) {
			t.Errorf("unexpected detection of custom keyword %s without config\nfull output:\n%s", kw, out)
		}
	}
}

// TestE2E_Config_PartialDisabled uses a config that disables first_letter, emoji,
// and security. The English filter remains enabled and must catch Cyrillic text;
// uppercase first letters and sensitive variable names must not be reported.
func TestE2E_Config_PartialDisabled(t *testing.T) {
	binary := buildLingo(t)
	projectDir := filepath.Join("testdata", "partial-disabled-project")
	configPath := absConfig(t, "partial-disabled-project")

	out, err := runVetWithConfig(t, binary, projectDir, configPath)
	if err == nil {
		t.Fatalf("expected English-filter issues but lingo found nothing\noutput:\n%s", out)
	}

	// English filter must fire on Cyrillic text
	if !strings.Contains(out, "must be in English") {
		t.Errorf("expected English filter diagnostic\nfull output:\n%s", out)
	}

	// first_letter is disabled — no uppercase diagnostics expected
	if strings.Contains(out, "must start with a lowercase letter") {
		t.Errorf("first_letter filter should be disabled but reported an issue\nfull output:\n%s", out)
	}

	// security is disabled — sensitive variable names must not be reported
	if strings.Contains(out, "may expose sensitive data") {
		t.Errorf("security filter should be disabled but reported an issue\nfull output:\n%s", out)
	}

	// emoji filter is disabled — no emoji diagnostics expected
	if strings.Contains(out, "must not contain emoji") {
		t.Errorf("emoji filter should be disabled but reported an issue\nfull output:\n%s", out)
	}
}

// TestE2E_Config_InvalidConfig verifies that a malformed JSON config causes
// a readable linter error instead of a panic or silent failure.
func TestE2E_Config_InvalidConfig(t *testing.T) {
	binary := buildLingo(t)
	projectDir := filepath.Join("testdata", "invalid-config-project")
	configPath := absConfig(t, "invalid-config-project")

	out, err := runVetWithConfig(t, binary, projectDir, configPath)
	if err == nil {
		t.Fatalf("expected error for invalid config, but go vet exited successfully\noutput:\n%s", out)
	}

	// the error output must mention the linter name
	if !strings.Contains(out, "lingo") {
		t.Errorf("expected error message to mention 'lingo'\nfull output:\n%s", out)
	}
}

// TestE2E_Config_NonExistentConfig verifies that a path to a missing config
// file produces a readable error.
func TestE2E_Config_NonExistentConfig(t *testing.T) {
	binary := buildLingo(t)
	projectDir := filepath.Join("testdata", "clean-project")
	configPath := "/nonexistent/path/.lingo.json"

	out, err := runVetWithConfig(t, binary, projectDir, configPath)
	if err == nil {
		t.Fatalf("expected error for missing config file\noutput:\n%s", out)
	}

	if !strings.Contains(out, "lingo") {
		t.Errorf("expected error message to mention 'lingo'\nfull output:\n%s", out)
	}
}
