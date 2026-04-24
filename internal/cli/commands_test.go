package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func tempCommander(t *testing.T) (*Commander, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "store.json")
	cmd, err := NewCommander(path)
	if err != nil {
		t.Fatalf("NewCommander: %v", err)
	}
	return cmd, path
}

func TestCommander_SetAndList(t *testing.T) {
	cmd, _ := tempCommander(t)

	if err := cmd.SetVar("dev", "DB_HOST", "localhost"); err != nil {
		t.Fatalf("SetVar: %v", err)
	}
	if err := cmd.SetVar("dev", "DB_PORT", "5432"); err != nil {
		t.Fatalf("SetVar: %v", err)
	}

	vars, err := cmd.resolver.Resolve([]string{"dev"})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if vars["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", vars["DB_HOST"])
	}
	if vars["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", vars["DB_PORT"])
	}
}

func TestCommander_DeleteVar(t *testing.T) {
	cmd, _ := tempCommander(t)

	_ = cmd.SetVar("prod", "SECRET", "abc123")
	if err := cmd.DeleteVar("prod", "SECRET"); err != nil {
		t.Fatalf("DeleteVar: %v", err)
	}

	vars, _ := cmd.resolver.Resolve([]string{"prod"})
	if _, ok := vars["SECRET"]; ok {
		t.Error("expected SECRET to be deleted")
	}
}

func TestCommander_DeleteVar_MissingProfile(t *testing.T) {
	cmd, _ := tempCommander(t)
	err := cmd.DeleteVar("ghost", "KEY")
	if err == nil {
		t.Error("expected error for missing profile")
	}
}

func TestCommander_ListProfiles(t *testing.T) {
	cmd, _ := tempCommander(t)
	_ = cmd.SetVar("alpha", "X", "1")
	_ = cmd.SetVar("beta", "Y", "2")

	// Capture stdout via redirect
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := cmd.ListProfiles()
	w.Close()
	os.Stdout = old

	buf := make([]byte, 256)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	if err != nil {
		t.Fatalf("ListProfiles: %v", err)
	}
	if output == "" {
		t.Error("expected profile listing output")
	}
}
