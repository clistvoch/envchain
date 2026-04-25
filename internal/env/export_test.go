package env_test

import (
	"strings"
	"testing"

	"github.com/user/envchain/internal/env"
	"github.com/user/envchain/internal/store"
)

func setupExporter(t *testing.T) *env.Exporter {
	t.Helper()
	s := store.NewStore()

	p1, _ := store.NewProfile("base")
	p1.Set("APP_ENV", "production")
	p1.Set("LOG_LEVEL", "info")

	p2, _ := store.NewProfile("override")
	p2.Set("LOG_LEVEL", "debug")
	p2.Set("SECRET", "my secret'value")

	_ = s.AddProfile(p1)
	_ = s.AddProfile(p2)

	r := env.NewResolver(s)
	return env.NewExporter(r)
}

func TestExporter_Bash(t *testing.T) {
	ex := setupExporter(t)
	var buf strings.Builder
	if err := ex.Export(&buf, env.FormatBash, "base"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export APP_ENV='production'") {
		t.Errorf("expected bash export for APP_ENV, got:\n%s", out)
	}
	if !strings.Contains(out, "export LOG_LEVEL='info'") {
		t.Errorf("expected bash export for LOG_LEVEL, got:\n%s", out)
	}
}

func TestExporter_Fish(t *testing.T) {
	ex := setupExporter(t)
	var buf strings.Builder
	if err := ex.Export(&buf, env.FormatFish, "base"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "set -x APP_ENV 'production'") {
		t.Errorf("expected fish set for APP_ENV, got:\n%s", out)
	}
}

func TestExporter_Dotenv(t *testing.T) {
	ex := setupExporter(t)
	var buf strings.Builder
	if err := ex.Export(&buf, env.FormatDotenv, "base"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected dotenv line for APP_ENV, got:\n%s", out)
	}
}

func TestExporter_ChainOverride(t *testing.T) {
	ex := setupExporter(t)
	var buf strings.Builder
	if err := ex.Export(&buf, env.FormatBash, "base", "override"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export LOG_LEVEL='debug'") {
		t.Errorf("expected overridden LOG_LEVEL=debug, got:\n%s", out)
	}
}

func TestExporter_ShellEscapeQuotes(t *testing.T) {
	ex := setupExporter(t)
	var buf strings.Builder
	if err := ex.Export(&buf, env.FormatBash, "override"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `'my secret'\''value'`) {
		t.Errorf("expected escaped single quote in SECRET, got:\n%s", out)
	}
}

func TestExporter_MissingProfile(t *testing.T) {
	ex := setupExporter(t)
	var buf strings.Builder
	err := ex.Export(&buf, env.FormatBash, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing profile, got nil")
	}
}
