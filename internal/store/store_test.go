package store

import (
	"os"
	"path/filepath"
	"testing"
)

func tempStorePath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "envchain.json")
}

func TestStore_AddAndGetProfile(t *testing.T) {
	s, err := NewStore(tempStorePath(t))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	p, _ := NewProfile("dev")
	p.Set("DB_URL", "postgres://localhost/dev")

	if err := s.AddProfile(p); err != nil {
		t.Fatalf("AddProfile: %v", err)
	}

	got, err := s.GetProfile("dev")
	if err != nil {
		t.Fatalf("GetProfile: %v", err)
	}
	if v, _ := got.Get("DB_URL"); v != "postgres://localhost/dev" {
		t.Errorf("unexpected DB_URL: %q", v)
	}
}

func TestStore_AddProfile_Duplicate(t *testing.T) {
	s, _ := NewStore(tempStorePath(t))
	p, _ := NewProfile("dev")
	s.AddProfile(p)
	if err := s.AddProfile(p); err != ErrProfileAlreadyExists {
		t.Errorf("expected ErrProfileAlreadyExists, got %v", err)
	}
}

func TestStore_DeleteProfile(t *testing.T) {
	s, _ := NewStore(tempStorePath(t))
	p, _ := NewProfile("staging")
	s.AddProfile(p)
	if err := s.DeleteProfile("staging"); err != nil {
		t.Fatalf("DeleteProfile: %v", err)
	}
	if _, err := s.GetProfile("staging"); err != ErrProfileNotFound {
		t.Errorf("expected ErrProfileNotFound after delete, got %v", err)
	}
}

func TestStore_SaveAndLoad(t *testing.T) {
	path := tempStorePath(t)
	s, _ := NewStore(path)
	p, _ := NewProfile("prod")
	p.Set("API_KEY", "secret")
	s.AddProfile(p)
	if err := s.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Reload from disk
	s2, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore reload: %v", err)
	}
	got, err := s2.GetProfile("prod")
	if err != nil {
		t.Fatalf("GetProfile after reload: %v", err)
	}
	if v, _ := got.Get("API_KEY"); v != "secret" {
		t.Errorf("expected API_KEY=secret after reload, got %q", v)
	}
}

func TestStore_FilePermissions(t *testing.T) {
	path := tempStorePath(t)
	s, _ := NewStore(path)
	p, _ := NewProfile("secure")
	s.AddProfile(p)
	s.Save()

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected file perm 0600, got %o", perm)
	}
}
