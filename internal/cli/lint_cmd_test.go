package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/urfave/cli/v2"
)

func setupLintApp(t *testing.T) (*cli.App, *Commander) {
	t.Helper()
	cmdr := tempCommander(t)
	app := &cli.App{Writer: &bytes.Buffer{}}
	RegisterLintCmd(app, cmdr)
	return app, cmdr
}

func TestLintCmd_NoIssues(t *testing.T) {
	app, cmdr := setupLintApp(t)

	_ = cmdr.Set("base", "DB_HOST", "localhost")
	_ = cmdr.Set("base", "DB_PORT", "5432")

	err := app.Run([]string{"app", "lint", "base"})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	out := app.Writer.(*bytes.Buffer).String()
	if !strings.Contains(out, "OK") {
		t.Errorf("expected OK message, got: %s", out)
	}
}

func TestLintCmd_ShadowedKey(t *testing.T) {
	app, cmdr := setupLintApp(t)

	_ = cmdr.Set("base", "DB_HOST", "localhost")
	_ = cmdr.Set("override", "DB_HOST", "prod-host")

	buf := &bytes.Buffer{}
	app.Writer = buf

	err := app.Run([]string{"app", "lint", "base", "override"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "shadows") {
		t.Errorf("expected shadow warning, got: %s", out)
	}
}

func TestLintCmd_InvalidVarName(t *testing.T) {
	app, cmdr := setupLintApp(t)

	// Manually insert an invalid key by bypassing Set validation
	p, _ := cmdr.store.GetProfile("bad")
	if p == nil {
		_ = cmdr.Set("bad", "VALID_KEY", "val")
		p, _ = cmdr.store.GetProfile("bad")
	}
	p.Vars["123INVALID"] = "oops"
	_ = cmdr.store.Save()

	buf := &bytes.Buffer{}
	app.Writer = buf

	_ = app.Run([]string{"app", "lint", "bad"})
	out := buf.String()
	if !strings.Contains(out, "invalid variable name") {
		t.Errorf("expected invalid name warning, got: %s", out)
	}
}

func TestLintCmd_MissingProfile(t *testing.T) {
	app, _ := setupLintApp(t)
	buf := &bytes.Buffer{}
	app.Writer = buf

	_ = app.Run([]string{"app", "lint", "nonexistent"})
	out := buf.String()
	if !strings.Contains(out, "not found") {
		t.Errorf("expected not-found warning, got: %s", out)
	}
}

func TestLintCmd_StrictMode(t *testing.T) {
	app, cmdr := setupLintApp(t)

	_ = cmdr.Set("a", "KEY", "1")
	_ = cmdr.Set("b", "KEY", "2")

	err := app.Run([]string{"app", "lint", "--strict", "a", "b"})
	if err == nil {
		t.Fatal("expected error in strict mode when warnings exist")
	}
}

func TestLintCmd_NoArgs(t *testing.T) {
	app, _ := setupLintApp(t)
	err := app.Run([]string{"app", "lint"})
	if err == nil {
		t.Fatal("expected error when no profiles provided")
	}
}

func TestIsValidEnvKey(t *testing.T) {
	cases := []struct {
		key   string
		valid bool
	}{
		{"MY_VAR", true},
		{"myvar", true},
		{"_PRIVATE", true},
		{"VAR123", true},
		{"123VAR", false},
		{"", false},
		{"MY-VAR", false},
		{"__DOUBLE", false},
	}
	for _, tc := range cases {
		got := isValidEnvKey(tc.key)
		if got != tc.valid {
			t.Errorf("isValidEnvKey(%q) = %v, want %v", tc.key, got, tc.valid)
		}
	}
}
