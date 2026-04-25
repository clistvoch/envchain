package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestDiffCmd_Identical(t *testing.T) {
	cmdr, _ := tempCommander(t)

	_ = cmdr.SetVar("alpha", "KEY", "value1")
	_ = cmdr.SetVar("beta", "KEY", "value1")

	var buf bytes.Buffer
	cmdr.out = &buf

	if err := cmdr.DiffProfiles("alpha", "beta"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "identical") {
		t.Errorf("expected 'identical' in output, got: %s", buf.String())
	}
}

func TestDiffCmd_OnlyInA(t *testing.T) {
	cmdr, _ := tempCommander(t)

	_ = cmdr.SetVar("alpha", "ONLY_A", "yes")
	_ = cmdr.SetVar("beta", "OTHER", "no")

	var buf bytes.Buffer
	cmdr.out = &buf

	if err := cmdr.DiffProfiles("alpha", "beta"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "- ONLY_A") {
		t.Errorf("expected '- ONLY_A' in output, got: %s", out)
	}
	if !strings.Contains(out, "+ OTHER") {
		t.Errorf("expected '+ OTHER' in output, got: %s", out)
	}
}

func TestDiffCmd_ChangedValue(t *testing.T) {
	cmdr, _ := tempCommander(t)

	_ = cmdr.SetVar("alpha", "SHARED", "old")
	_ = cmdr.SetVar("beta", "SHARED", "new")

	var buf bytes.Buffer
	cmdr.out = &buf

	if err := cmdr.DiffProfiles("alpha", "beta"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "~ SHARED") {
		t.Errorf("expected '~ SHARED' in output, got: %s", out)
	}
}

func TestDiffCmd_MissingProfile(t *testing.T) {
	cmdr, _ := tempCommander(t)

	_ = cmdr.SetVar("alpha", "KEY", "val")

	if err := cmdr.DiffProfiles("alpha", "ghost"); err == nil {
		t.Error("expected error for missing profile")
	}
}
