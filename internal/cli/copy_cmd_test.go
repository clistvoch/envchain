package cli

import (
	"bytes"
	"testing"
)

func TestCopyCmd_Basic(t *testing.T) {
	c := tempCommander(t)

	// Populate source profile.
	if err := c.store.AddProfile("src"); err != nil {
		t.Fatal(err)
	}
	src, _ := c.store.GetProfile("src")
	src.Set("FOO", "bar")
	src.Set("BAZ", "qux")
	if err := c.store.Save(); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	cmd := c.newCopyCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"src", "dst"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dst, err := c.store.GetProfile("dst")
	if err != nil {
		t.Fatalf("destination profile not found: %v", err)
	}
	if v, ok := dst.Get("FOO"); !ok || v != "bar" {
		t.Errorf("expected FOO=bar, got %q (ok=%v)", v, ok)
	}
	if v, ok := dst.Get("BAZ"); !ok || v != "qux" {
		t.Errorf("expected BAZ=qux, got %q (ok=%v)", v, ok)
	}
}

func TestCopyCmd_NoOverwriteByDefault(t *testing.T) {
	c := tempCommander(t)

	for _, name := range []string{"src", "dst"} {
		if err := c.store.AddProfile(name); err != nil {
			t.Fatal(err)
		}
	}
	src, _ := c.store.GetProfile("src")
	src.Set("KEY", "new")
	dst, _ := c.store.GetProfile("dst")
	dst.Set("KEY", "original")
	if err := c.store.Save(); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	cmd := c.newCopyCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"src", "dst"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dst, _ = c.store.GetProfile("dst")
	if v, _ := dst.Get("KEY"); v != "original" {
		t.Errorf("expected KEY to remain 'original', got %q", v)
	}
}

func TestCopyCmd_Overwrite(t *testing.T) {
	c := tempCommander(t)

	for _, name := range []string{"src", "dst"} {
		if err := c.store.AddProfile(name); err != nil {
			t.Fatal(err)
		}
	}
	src, _ := c.store.GetProfile("src")
	src.Set("KEY", "new")
	dst, _ := c.store.GetProfile("dst")
	dst.Set("KEY", "original")
	if err := c.store.Save(); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	cmd := c.newCopyCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--overwrite", "src", "dst"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dst, _ = c.store.GetProfile("dst")
	if v, _ := dst.Get("KEY"); v != "new" {
		t.Errorf("expected KEY='new' after overwrite, got %q", v)
	}
}

func TestCopyCmd_MissingSource(t *testing.T) {
	c := tempCommander(t)
	cmd := c.newCopyCmd()
	cmd.SetArgs([]string{"nonexistent", "dst"})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error for missing source profile, got nil")
	}
}
