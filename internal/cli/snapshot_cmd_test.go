package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestSnapshotSave_Basic(t *testing.T) {
	app, cmdr := tempCommander(t)

	if err := cmdr.SetVar("prod", "DB_URL", "postgres://localhost/prod"); err != nil {
		t.Fatal(err)
	}

	err := app.Run([]string{"envchain", "snapshot", "save", "prod", "v1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	profiles, _ := cmdr.store.ListProfiles()
	found := false
	for _, p := range profiles {
		if p == "prod.snap.v1" {
			found = true
		}
	}
	if !found {
		t.Error("expected snapshot profile prod.snap.v1 to exist")
	}
}

func TestSnapshotList_ShowsSnapshots(t *testing.T) {
	app, cmdr := tempCommander(t)

	if err := cmdr.SetVar("staging", "API_KEY", "abc"); err != nil {
		t.Fatal(err)
	}
	_ = app.Run([]string{"envchain", "snapshot", "save", "staging", "before-deploy"})
	_ = app.Run([]string{"envchain", "snapshot", "save", "staging", "after-deploy"})

	var buf bytes.Buffer
	app.Writer = &buf
	err := app.Run([]string{"envchain", "snapshot", "list", "staging"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "before-deploy") {
		t.Errorf("expected 'before-deploy' in output, got: %s", out)
	}
	if !strings.Contains(out, "after-deploy") {
		t.Errorf("expected 'after-deploy' in output, got: %s", out)
	}
}

func TestSnapshotRestore_Basic(t *testing.T) {
	app, cmdr := tempCommander(t)

	if err := cmdr.SetVar("dev", "SECRET", "original"); err != nil {
		t.Fatal(err)
	}
	_ = app.Run([]string{"envchain", "snapshot", "save", "dev", "checkpoint"})

	if err := cmdr.SetVar("dev", "SECRET", "changed"); err != nil {
		t.Fatal(err)
	}

	err := app.Run([]string{"envchain", "snapshot", "restore", "--overwrite", "dev", "checkpoint"})
	if err != nil {
		t.Fatalf("unexpected restore error: %v", err)
	}

	val, err := cmdr.GetVar("dev", "SECRET")
	if err != nil {
		t.Fatal(err)
	}
	if val != "original" {
		t.Errorf("expected 'original' after restore, got %q", val)
	}
}

func TestSnapshotList_Empty(t *testing.T) {
	app, cmdr := tempCommander(t)
	if err := cmdr.SetVar("empty", "X", "1"); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	app.Writer = &buf
	err := app.Run([]string{"envchain", "snapshot", "list", "empty"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No snapshots") {
		t.Errorf("expected 'No snapshots' message, got: %s", buf.String())
	}
}
