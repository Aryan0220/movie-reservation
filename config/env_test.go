package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetEnv_ReturnsValue(t *testing.T) {
	t.Setenv("TEST_KEY", "test-value")

	if got := GetEnv("TEST_KEY"); got != "test-value" {
		t.Fatalf("expected test-value, got %q", got)
	}
}

func TestLoadEnv_FindsDotEnv(t *testing.T) {
	tempDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldDir)
	}()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change dir: %v", err)
	}

	envPath := filepath.Join(tempDir, ".env")
	if err := os.WriteFile(envPath, []byte("PORT=1234\n"), 0600); err != nil {
		t.Fatalf("failed to write .env: %v", err)
	}

	if got := LoadEnv(); got != ".Env File Found" {
		t.Fatalf("expected .Env File Found, got %q", got)
	}
}
