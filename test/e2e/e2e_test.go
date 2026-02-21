//go:build e2e

// Package e2e тестирует полный цикл работы линтера:
// сборка бинарника cmd/addcheck → запуск через go vet → проверка диагностик.
//
// Запуск: go test -tags e2e ./test/e2e/
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
	// test/e2e/ → на два уровня вверх
	return filepath.Join(wd, "..", "..")
}

func buildLingo(t *testing.T) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "lingo-check")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	root := projectRoot(t)
	cmd := exec.Command("go", "build", "-o", binary, filepath.Join(root, "cmd", "addcheck"))
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build lingo: %s\n%s", err, out)
	}
	return binary
}

// runVet запускает go vet без конфига в указанной директории.
func runVet(t *testing.T, binary, projectDir string) (string, error) {
	t.Helper()
	cmd := exec.Command("go", "vet", "-vettool="+binary, "./...")
	cmd.Dir = projectDir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// runVetWithConfig запускает go vet с флагом -config.
// configPath должен быть абсолютным путём к .lingo.json.
func runVetWithConfig(t *testing.T, binary, projectDir, configPath string) (string, error) {
	t.Helper()
	cmd := exec.Command("go", "vet", "-vettool="+binary, "-config="+configPath, "./...")
	cmd.Dir = projectDir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// absConfig возвращает абсолютный путь к .lingo.json внутри testdata-проекта.
func absConfig(t *testing.T, projectName string) string {
	t.Helper()
	path, err := filepath.Abs(filepath.Join("testdata", projectName, ".lingo.json"))
	if err != nil {
		t.Fatal(err)
	}
	return path
}

// ── Существующие тесты ──────────────────────────────────────────────────────

// TestE2E_CleanProject — линтер не должен ничего сообщать на чистом коде.
func TestE2E_CleanProject(t *testing.T) {
	binary := buildLingo(t)
	projectDir := filepath.Join("testdata", "clean-project")

	out, err := runVet(t, binary, projectDir)
	if err != nil {
		t.Errorf("lingo reported issues on clean code (expected no issues):\n%s", out)
	}
}

// TestE2E_ViolationsProject — линтер должен найти все виды нарушений.
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

// ── Тесты с конфигурацией ───────────────────────────────────────────────────

// TestE2E_Config_DisabledFilters — если все фильтры отключены в конфиге,
// линтер не должен репортить ни одной ошибки, даже на коде с нарушениями.
func TestE2E_Config_DisabledFilters(t *testing.T) {
	binary := buildLingo(t)
	projectDir := filepath.Join("testdata", "disabled-filters-project")
	configPath := absConfig(t, "disabled-filters-project")

	out, err := runVetWithConfig(t, binary, projectDir, configPath)
	if err != nil {
		t.Errorf("expected no issues when all filters disabled, but got:\n%s", out)
	}
}

// TestE2E_Config_DisabledFilters_WithoutConfig — тот же код с нарушениями, но
// без конфига — линтер должен найти ошибки (фильтры включены по умолчанию).
func TestE2E_Config_DisabledFilters_WithoutConfig(t *testing.T) {
	binary := buildLingo(t)
	projectDir := filepath.Join("testdata", "disabled-filters-project")

	out, err := runVet(t, binary, projectDir)
	if err == nil {
		t.Fatalf("expected issues without config, but lingo reported nothing\noutput:\n%s", out)
	}
}

// TestE2E_Config_ExtraKeywords — кастомные keywords (cvv, ssn, otp) из конфига
// должны детектироваться наравне со встроенными (password).
func TestE2E_Config_ExtraKeywords(t *testing.T) {
	binary := buildLingo(t)
	projectDir := filepath.Join("testdata", "extra-keywords-project")
	configPath := absConfig(t, "extra-keywords-project")

	out, err := runVetWithConfig(t, binary, projectDir, configPath)
	if err == nil {
		t.Fatalf("expected security issues but lingo found nothing\noutput:\n%s", out)
	}

	// каждый кастомный keyword должен быть задетектирован
	expectedKeywords := []string{"cvv", "ssn", "otp"}
	for _, kw := range expectedKeywords {
		if !strings.Contains(out, kw) {
			t.Errorf("expected report for custom keyword %q not found\nfull output:\n%s", kw, out)
		}
	}

	// встроенный keyword тоже должен работать
	if !strings.Contains(out, "password") {
		t.Errorf("expected report for built-in keyword %q not found\nfull output:\n%s", "password", out)
	}

	// безопасная переменная requestID не должна срабатывать
	if strings.Contains(out, "requestID") {
		t.Errorf("false positive: requestID should not trigger security filter\nfull output:\n%s", out)
	}

	// все диагностики должны содержать стандартное сообщение
	if !strings.Contains(out, "may expose sensitive data") {
		t.Errorf("expected standard security message not found\nfull output:\n%s", out)
	}
}

// TestE2E_Config_ExtraKeywords_NotDetectedWithoutConfig — без конфига кастомные
// keywords (cvv, ssn, otp) не знакомы линтеру, "otp code sent" не срабатывает.
func TestE2E_Config_ExtraKeywords_NotDetectedWithoutConfig(t *testing.T) {
	binary := buildLingo(t)
	projectDir := filepath.Join("testdata", "extra-keywords-project")

	out, err := runVet(t, binary, projectDir)

	// встроенный "password" всё ещё должен сработать
	if err == nil {
		t.Fatalf("expected at least one issue (password), but lingo found nothing\noutput:\n%s", out)
	}
	if !strings.Contains(out, "password") {
		t.Errorf("expected built-in 'password' detection\nfull output:\n%s", out)
	}

	// кастомные keywords НЕ должны срабатывать без конфига
	for _, kw := range []string{"\"cvv\"", "\"ssn\"", "\"otp\""} {
		if strings.Contains(out, kw) {
			t.Errorf("unexpected detection of custom keyword %s without config\nfull output:\n%s", kw, out)
		}
	}
}

// TestE2E_Config_PartialDisabled — конфиг отключает first_letter, emoji, security.
// English-фильтр остаётся включённым и должен поймать кириллицу.
// Заглавные буквы при этом не должны вызывать диагностику.
func TestE2E_Config_PartialDisabled(t *testing.T) {
	binary := buildLingo(t)
	projectDir := filepath.Join("testdata", "partial-disabled-project")
	configPath := absConfig(t, "partial-disabled-project")

	out, err := runVetWithConfig(t, binary, projectDir, configPath)
	if err == nil {
		t.Fatalf("expected English-filter issues but lingo found nothing\noutput:\n%s", out)
	}

	// English-фильтр должен сработать на кириллице
	if !strings.Contains(out, "must be in English") {
		t.Errorf("expected English filter diagnostic\nfull output:\n%s", out)
	}

	// first_letter отключён — заглавных ошибок быть не должно
	if strings.Contains(out, "must start with a lowercase letter") {
		t.Errorf("first_letter filter should be disabled but reported an issue\nfull output:\n%s", out)
	}

	// security отключён — переменная token не должна срабатывать
	if strings.Contains(out, "may expose sensitive data") {
		t.Errorf("security filter should be disabled but reported an issue\nfull output:\n%s", out)
	}

	// emoji отключён — не должно быть диагностик по emoji
	if strings.Contains(out, "must not contain emoji") {
		t.Errorf("emoji filter should be disabled but reported an issue\nfull output:\n%s", out)
	}
}

// TestE2E_Config_InvalidConfig — невалидный JSON в конфиге должен приводить к
// ошибке линтера, а не panic или silent fail.
func TestE2E_Config_InvalidConfig(t *testing.T) {
	binary := buildLingo(t)
	projectDir := filepath.Join("testdata", "invalid-config-project")
	configPath := absConfig(t, "invalid-config-project")

	out, err := runVetWithConfig(t, binary, projectDir, configPath)
	if err == nil {
		t.Fatalf("expected error for invalid config, but go vet exited successfully\noutput:\n%s", out)
	}

	// ошибка должна содержать внятное сообщение о конфиге
	if !strings.Contains(out, "lingo") {
		t.Errorf("expected error message to mention 'lingo'\nfull output:\n%s", out)
	}
}

// TestE2E_Config_NonExistentConfig — путь к несуществующему конфигу должен
// вернуть читаемую ошибку.
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
