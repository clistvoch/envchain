package cli

import (
	"bytes"
	"testing"
)

func TestRenameCmd_Basic(t *testing.T) {
	dir := t.TempDir()
	cmdr := tempCommander(t, dir)

	// Populate source profile
	if err := cmdr.SetVar("dev", "KEY", "value1"); err != nil {
		t.Fatalf("SetVar: %v", err)
	}
	if err := cmdr.SetVar("dev", "KEY2", "value2"); err != nil {
		t.Fatalf("SetVar: %v", err)
	}

	if err := cmdr.RenameProfile("dev", "staging"); err != nil {
		t.Fatalf("RenameProfile: %v", err)
	}

	// Old profile must be gone
	if _, err := cmdr.store.GetProfile("dev"); err == nil {
		t.Error("expected old profile 'dev' to be deleted")
	}

	// New profile must carry all vars
	p, err := cmdr.store.GetProfile("staging")
	if err != nil {
		t.Fatalf("GetProfile staging: %v", err)
	}
	if v, _ := p.Get("KEY"); v != "value1" {
		t.Errorf("expected KEY=value1, got %q", v)
	}
	if v, _ := p.Get("KEY2"); v != "value2" {
		t.Errorf("expected KEY2=value2, got %q", v)
	}
}

func TestRenameCmd_MissingSource(t *testing.T) {
	dir := t.TempDir()
	cmdr := tempCommander(t, dir)

	err := cmdr.RenameProfile("nonexistent", "other")
	if err == nil {
		t.Fatal("expected error for missing source profile")
	}
}

func TestRenameCmd_DestinationExists(t *testing.T) {
	dir := t.TempDir()
	cmdr := tempCommander(t, dir)

	if err := cmdr.SetVar("alpha", "X", "1"); err != nil {
		t.Fatalf("SetVar: %v", err)
	}
	if err := cmdr.SetVar("beta", "Y", "2"); err != nil {
		t.Fatalf("SetVar: %v", err)
	}

	err := cmdr.RenameProfile("alpha", "beta")
	if err == nil {
		t.Fatal("expected error when destination profile already exists")
	}
}

func TestRenameCmd_SameName(t *testing.T) {
	dir := t.TempDir()
	cmdr := tempCommander(t, dir)

	err := cmdr.RenameProfile("dev", "dev")
	if err == nil {
		t.Fatal("expected error when old and new names are identical")
	}
}

func TestRenameCmd_OutputMessage(t *testing.T) {
	dir := t.TempDir()
	cmdr := tempCommander(t, dir)

	var buf bytes.Buffer
	cmdr.out = &buf

	if err := cmdr.SetVar("prod", "TOKEN", "abc"); err != nil {
		t.Fatalf("SetVar: %v", err)
	}
	if err := cmdr.RenameProfile("prod", "production"); err != nil {
		t.Fatalf("RenameProfile: %v", err)
	}

	out := buf.String()
	if out == "" {
		t.Error("expected output message after rename")
	}
}
