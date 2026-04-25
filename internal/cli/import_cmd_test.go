package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempDotenv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp env file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp env file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestImportDotenv_Basic(t *testing.T) {
	c := tempCommander(t)
	envFile := writeTempDotenv(t, "# comment\nFOO=bar\nBAZ=qux\n")

	if err := c.ImportDotenv("default", envFile); err != nil {
		t.Fatalf("ImportDotenv: %v", err)
	}

	vars, err := c.ListVars("default")
	if err != nil {
		t.Fatalf("ListVars: %v", err)
	}
	if vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", vars["FOO"])
	}
	if vars["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", vars["BAZ"])
	}
}

func TestImportDotenv_QuotedValues(t *testing.T) {
	c := tempCommander(t)
	envFile := writeTempDotenv(t, `KEY1="hello world"
KEY2='single quoted'
export KEY3="exported"
`)

	if err := c.ImportDotenv("default", envFile); err != nil {
		t.Fatalf("ImportDotenv: %v", err)
	}

	vars, _ := c.ListVars("default")
	if vars["KEY1"] != "hello world" {
		t.Errorf("KEY1: got %q", vars["KEY1"])
	}
	if vars["KEY2"] != "single quoted" {
		t.Errorf("KEY2: got %q", vars["KEY2"])
	}
	if vars["KEY3"] != "exported" {
		t.Errorf("KEY3: got %q", vars["KEY3"])
	}
}

func TestImportDotenv_MissingFile(t *testing.T) {
	c := tempCommander(t)
	err := c.ImportDotenv("default", filepath.Join(t.TempDir(), "nonexistent.env"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestImportDotenv_InvalidLine(t *testing.T) {
	c := tempCommander(t)
	envFile := writeTempDotenv(t, "NOEQUALSSIGN\n")
	err := c.ImportDotenv("default", envFile)
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}
