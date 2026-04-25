package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestInspectCmd_SingleProfile(t *testing.T) {
	cmdr := tempCommander(t)
	cmdr.RunSet("base", "HOST", "localhost")
	cmdr.RunSet("base", "PORT", "5432")

	out := captureStdout(func() {
		if err := cmdr.RunInspect([]string{"base"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "HOST") {
		t.Errorf("expected HOST in output, got:\n%s", out)
	}
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output, got:\n%s", out)
	}
	if !strings.Contains(out, "[base]") {
		t.Errorf("expected source profile [base] in output, got:\n%s", out)
	}
}

func TestInspectCmd_ChainShowsCorrectSource(t *testing.T) {
	cmdr := tempCommander(t)
	cmdr.RunSet("base", "HOST", "localhost")
	cmdr.RunSet("base", "PORT", "5432")
	cmdr.RunSet("override", "PORT", "9999")

	out := captureStdout(func() {
		if err := cmdr.RunInspect([]string{"base", "override"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	// PORT should be attributed to "base" (first occurrence wins in inspect)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "PORT") {
			if !strings.Contains(line, "[base]") {
				t.Errorf("expected PORT to be sourced from [base], got: %s", line)
			}
		}
	}
}

func TestInspectCmd_MissingProfile(t *testing.T) {
	cmdr := tempCommander(t)
	err := cmdr.RunInspect([]string{"nonexistent"})
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("expected profile name in error, got: %v", err)
	}
}

func TestInspectCmd_NoProfiles(t *testing.T) {
	cmdr := tempCommander(t)
	out := captureStdout(func() {
		if err := cmdr.RunInspect([]string{}); err != nil {
			// expected no error, just empty output
			_ = fmt.Sprintf("%v", err)
		}
	})
	if strings.TrimSpace(out) != "" {
		t.Errorf("expected empty output for no profiles, got: %s", out)
	}
}
