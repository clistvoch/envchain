package cli

import (
	"testing"
)

// TestLintProfiles_EmptyChain ensures linting an empty slice returns no warnings.
func TestLintProfiles_EmptyChain(t *testing.T) {
	cmdr := tempCommander(t)
	warnings := lintProfiles(cmdr, []string{})
	if len(warnings) != 0 {
		t.Errorf("expected 0 warnings for empty chain, got %d", len(warnings))
	}
}

// TestLintProfiles_MultipleShadows checks that each shadowed key produces its own warning.
func TestLintProfiles_MultipleShadows(t *testing.T) {
	cmdr := tempCommander(t)

	_ = cmdr.Set("p1", "A", "1")
	_ = cmdr.Set("p1", "B", "2")
	_ = cmdr.Set("p2", "A", "10")
	_ = cmdr.Set("p2", "B", "20")

	warnings := lintProfiles(cmdr, []string{"p1", "p2"})

	shadowCount := 0
	for _, w := range warnings {
		if len(w) > 0 {
			import_strings := w
			_ = import_strings
			shadowCount++
		}
	}

	if shadowCount < 2 {
		t.Errorf("expected at least 2 shadow warnings, got %d: %v", shadowCount, warnings)
	}
}

// TestLintProfiles_NoShadowSingleProfile ensures a single profile never has shadow warnings.
func TestLintProfiles_NoShadowSingleProfile(t *testing.T) {
	cmdr := tempCommander(t)

	_ = cmdr.Set("solo", "KEY_ONE", "v1")
	_ = cmdr.Set("solo", "KEY_TWO", "v2")

	warnings := lintProfiles(cmdr, []string{"solo"})
	for _, w := range warnings {
		if len(w) > 0 && contains(w, "shadows") {
			t.Errorf("unexpected shadow warning for single profile: %s", w)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 &&
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
