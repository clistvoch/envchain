package store

import (
	"testing"
)

func TestNewProfile_ValidName(t *testing.T) {
	p, err := NewProfile("my-profile_1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.Name != "my-profile_1" {
		t.Errorf("expected name 'my-profile_1', got %q", p.Name)
	}
}

func TestNewProfile_InvalidName(t *testing.T) {
	invalidNames := []string{"my profile", "prof@1", "", "prof/name"}
	for _, name := range invalidNames {
		_, err := NewProfile(name)
		if err == nil {
			t.Errorf("expected error for name %q, got nil", name)
		}
	}
}

func TestProfile_SetGet(t *testing.T) {
	p, _ := NewProfile("test")
	p.Set("FOO", "bar")
	v, ok := p.Get("FOO")
	if !ok || v != "bar" {
		t.Errorf("expected FOO=bar, got %q ok=%v", v, ok)
	}
}

func TestProfile_Delete(t *testing.T) {
	p, _ := NewProfile("test")
	p.Set("FOO", "bar")
	p.Delete("FOO")
	_, ok := p.Get("FOO")
	if ok {
		t.Error("expected FOO to be deleted")
	}
}

func TestProfile_Merge(t *testing.T) {
	base, _ := NewProfile("base")
	base.Set("A", "1")
	base.Set("B", "2")

	override, _ := NewProfile("override")
	override.Set("B", "99")
	override.Set("C", "3")

	base.Merge(override)

	if v, _ := base.Get("A"); v != "1" {
		t.Errorf("expected A=1, got %q", v)
	}
	if v, _ := base.Get("B"); v != "99" {
		t.Errorf("expected B=99 after merge, got %q", v)
	}
	if v, _ := base.Get("C"); v != "3" {
		t.Errorf("expected C=3, got %q", v)
	}
}
