package cli

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// RunWithProfiles resolves the given profile chain and executes the
// provided command with the merged environment injected.
func (c *Commander) RunWithProfiles(profiles []string, command string, args []string) error {
	if command == "" {
		return fmt.Errorf("command must not be empty")
	}

	vars, err := c.resolver.Resolve(profiles)
	if err != nil {
		return fmt.Errorf("resolving profiles: %w", err)
	}

	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start with the current process environment.
	cmd.Env = os.Environ()
	for k, v := range vars {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}
		return fmt.Errorf("command failed: %w", err)
	}
	return nil
}
