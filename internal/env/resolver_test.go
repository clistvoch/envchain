package env_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain/envchain/internal/env"
	"github.com/envchain/envchain/internal/store"
)

func setupStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	s, err := store.NewStore(filepath.Join(dir, "envchain.json"))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestResolver_Resolve_SingleProfile(t *testing.T) {
	s := setupStore(t)
	p, _ := store.NewProfile("base")
	p.Set("FOO", "bar")
	p.Set("BAZ", "qux")
	_ = s.AddProfile(p)

	r := env.NewResolver(s)
	vars, err := r.Resolve([]string{"base"})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if vars["FOO"] != "bar" || vars["BAZ"] != "qux" {
		t.Errorf("unexpected vars: %v", vars)
	}
}

func TestResolver_Resolve_ChainOverride(t *testing.T) {
	s := setupStore(t)

	base, _ := store.NewProfile("base")
	base.Set("FOO", "base_foo")
	base.Set("SHARED", "from_base")
	_ = s.AddProfile(base)

	override, _ := store.NewProfile("override")
	override.Set("SHARED", "from_override")
	override.Set("BAR", "baz")
	_ = s.AddProfile(override)

	r := env.NewResolver(s)
	vars, err := r.Resolve([]string{"base", "override"})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if vars["SHARED"] != "from_override" {
		t.Errorf("expected override, got %q", vars["SHARED"])
	}
	if vars["FOO"] != "base_foo" {
		t.Errorf("expected base_foo, got %q", vars["FOO"])
	}
}

func TestResolver_Resolve_MissingProfile(t *testing.T) {
	s := setupStore(t)
	r := env.NewResolver(s)
	_, err := r.Resolve([]string{"nonexistent"})
	if err == nil {
		t.Error("expected error for missing profile")
	}
}

func TestResolver_ApplyToProcess(t *testing.T) {
	s := setupStore(t)
	p, _ := store.NewProfile("ci")
	p.Set("ENVCHAIN_TEST_VAR", "hello")
	_ = s.AddProfile(p)

	r := env.NewResolver(s)
	if err := r.ApplyToProcess([]string{"ci"}); err != nil {
		t.Fatalf("ApplyToProcess: %v", err)
	}
	if got := os.Getenv("ENVCHAIN_TEST_VAR"); got != "hello" {
		t.Errorf("expected hello, got %q", got)
	}
	os.Unsetenv("ENVCHAIN_TEST_VAR")
}

func TestExpandVars(t *testing.T) {
	vars := map[string]string{"HOME": "/home/user", "APP": "envchain"}
	got := env.ExpandVars(vars, "${HOME}/.config/${APP}")
	want := "/home/user/.config/envchain"
	if got != want {
		t.Errorf("ExpandVars: got %q, want %q", got, want)
	}
}
