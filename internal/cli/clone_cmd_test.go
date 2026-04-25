package cli

import (
	"testing"
)

func TestCloneCmd_Basic(t *testing.T) {
	cmdr := tempCommander(t)

	if err := cmdr.SetVar("prod", "DB_HOST", "db.prod"); err != nil {
		t.Fatal(err)
	}
	if err := cmdr.SetVar("prod", "DB_PORT", "5432"); err != nil {
		t.Fatal(err)
	}

	if err := cmdr.CloneProfile("prod", "staging", "", false); err != nil {
		t.Fatalf("CloneProfile failed: %v", err)
	}

	v, err := cmdr.GetVar("staging", "DB_HOST")
	if err != nil {
		t.Fatalf("GetVar failed: %v", err)
	}
	if v != "db.prod" {
		t.Errorf("expected db.prod, got %q", v)
	}
}

func TestCloneCmd_WithPrefix(t *testing.T) {
	cmdr := tempCommander(t)

	_ = cmdr.SetVar("base", "HOST", "localhost")
	_ = cmdr.SetVar("base", "PORT", "8080")

	if err := cmdr.CloneProfile("base", "prefixed", "APP_", false); err != nil {
		t.Fatalf("CloneProfile with prefix failed: %v", err)
	}

	v, err := cmdr.GetVar("prefixed", "APP_HOST")
	if err != nil {
		t.Fatalf("expected APP_HOST in cloned profile: %v", err)
	}
	if v != "localhost" {
		t.Errorf("expected localhost, got %q", v)
	}

	if _, err := cmdr.GetVar("prefixed", "HOST"); err == nil {
		t.Error("expected original key HOST to be absent in prefixed clone")
	}
}

func TestCloneCmd_NoOverwriteByDefault(t *testing.T) {
	cmdr := tempCommander(t)

	_ = cmdr.SetVar("src", "KEY", "val")
	_ = cmdr.SetVar("dst", "KEY", "other")

	err := cmdr.CloneProfile("src", "dst", "", false)
	if err == nil {
		t.Fatal("expected error when destination exists and overwrite=false")
	}
}

func TestCloneCmd_Overwrite(t *testing.T) {
	cmdr := tempCommander(t)

	_ = cmdr.SetVar("src", "KEY", "new_val")
	_ = cmdr.SetVar("dst", "KEY", "old_val")

	if err := cmdr.CloneProfile("src", "dst", "", true); err != nil {
		t.Fatalf("CloneProfile overwrite failed: %v", err)
	}

	v, _ := cmdr.GetVar("dst", "KEY")
	if v != "new_val" {
		t.Errorf("expected new_val after overwrite, got %q", v)
	}
}

func TestCloneCmd_MissingSource(t *testing.T) {
	cmdr := tempCommander(t)

	err := cmdr.CloneProfile("nonexistent", "dst", "", false)
	if err == nil {
		t.Fatal("expected error for missing source profile")
	}
}

func TestCloneCmd_SameName(t *testing.T) {
	cmdr := tempCommander(t)

	_ = cmdr.SetVar("prod", "KEY", "val")

	err := cmdr.CloneProfile("prod", "prod", "", false)
	if err == nil {
		t.Fatal("expected error when source and destination are the same")
	}
}
