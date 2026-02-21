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

// TestE2E_CleanProject — линтер не должен ничего сообщать на чистом коде.
func TestE2E_CleanProject(t *testing.T) {
	binary := buildLingo(t)
	projectDir := filepath.Join("testdata", "clean-project")

	cmd := exec.Command("go", "vet", "-vettool="+binary, "./...")
	cmd.Dir = projectDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("lingo reported issues on clean code (expected no issues):\n%s", out)
	}
}

// TestE2E_ViolationsProject — линтер должен найти все виды нарушений.
func TestE2E_ViolationsProject(t *testing.T) {
	binary := buildLingo(t)
	projectDir := filepath.Join("testdata", "violations-project")

	cmd := exec.Command("go", "vet", "-vettool="+binary, "./...")
	cmd.Dir = projectDir
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("lingo found no issues in violations project (expected several)")
	}

	output := string(out)
	expected := []string{
		"must start with a lowercase letter",
		"must be in English",
		"must not contain emoji",
		"must not contain repeated punctuation",
		"may expose sensitive data",
	}
	for _, msg := range expected {
		if !strings.Contains(output, msg) {
			t.Errorf("expected diagnostic not found: %q\nfull output:\n%s", msg, output)
		}
	}
}
